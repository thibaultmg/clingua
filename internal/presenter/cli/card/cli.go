package card

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/looplab/fsm"
	"github.com/rs/zerolog/log"

	"github.com/thibaultmg/clingua/internal/entity"
)

type state string

const (
	initState              state = "init"
	editFieldState         state = "editField"
	editCardState          state = "editCard"
	saveCardState          state = "saveCard"
	endState               state = "end"
	selectPropositionState state = "selectProposition"
	writeFieldState        state = "writeField"
	listCardsState         state = "listCards"
	selectExampleState     state = "selectExample"
	editExampleState       state = "editExample"
)

type event string

const (
	editFieldEvent         event = "editField"
	editCardEvent          event = "editCard"
	selectPropositionEvent event = "selectProposition"
	writeFieldEvent        event = "writeField"
	listCardsEvent         event = "listCards"
	quitEvent              event = "quit"
	saveCardEvent          event = "saveCard"
	selectExampleEvent     event = "selectExample"
	editExampleEvent       event = "editExample"
)

type CardCLI struct {
	ce               *CardEditor
	FSM              *fsm.FSM
	console          console
	activeField      CardField
	activeFieldIndex int
	initTrack        bool
	done             <-chan struct{}
}

func NewCardCLI(ce *CardEditor) *CardCLI {
	doneChan := make(chan struct{})
	ret := &CardCLI{
		ce:      ce,
		done:    doneChan,
		console: newConsole(os.Stdout),
	}

	ret.FSM = fsm.NewFSM(
		string(initState),
		fsm.Events{
			{
				Name: string(editFieldEvent),
				Src:  []string{string(writeFieldState), string(selectPropositionState), string(editCardState), string(initState)},
				Dst:  string(editFieldState),
			},
			{
				Name: string(editCardEvent),
				Src:  []string{string(listCardsState), string(editFieldState), string(saveCardState)},
				Dst:  string(editCardState),
			},
			{
				Name: string(selectPropositionEvent),
				Src:  []string{string(editFieldState), string(editExampleState)},
				Dst:  string(selectPropositionState),
			},
			{Name: string(writeFieldEvent), Src: []string{string(editFieldState), string(editExampleState)}, Dst: string(writeFieldState)},
			{Name: string(quitEvent), Src: []string{string(editCardState)}, Dst: string(endState)},
			{Name: string(listCardsEvent), Src: []string{string(editCardState), string(initState)}, Dst: string(listCardsState)},
			{Name: string(saveCardEvent), Src: []string{string(editCardState)}, Dst: string(saveCardState)},
			{Name: string(selectExampleEvent), Src: []string{string(editCardState)}, Dst: string(selectExampleState)},
			{Name: string(editExampleEvent), Src: []string{string(selectExampleState), string(writeFieldState)}, Dst: string(editExampleState)},
		},
		fsm.Callbacks{
			"enter_state": func(e *fsm.Event) {
				log.Debug().Msgf("Enter state with event %s from state %s to state %s", e.Event, e.Src, e.Dst)
			},
			"enter_" + string(editFieldState):         ret.editFieldMenu,
			"enter_" + string(selectPropositionState): ret.showFieldPropositions,
			"enter_" + string(writeFieldState):        ret.editFieldPrompt,
			"enter_" + string(editCardState):          ret.cardMenu,
			"enter_" + string(saveCardState):          ret.save,
			"enter_" + string(listCardsState):         ret.listCards,
			"enter_" + string(selectExampleState):     ret.selectExampleMenu,
			"enter_" + string(editExampleState):       ret.editExampleMenu,
			"enter_" + string(endState): func(e *fsm.Event) {
				close(doneChan)
			},
		},
	)

	log.Debug().Msg(fsm.Visualize(ret.FSM))

	return ret
}

func (c *CardCLI) RunCreate() {
	c.initTrack = true
	c.activeField = DefinitionField
	c.sendEvent(editFieldEvent)
	<-c.done
}

func (c *CardCLI) RunList() {
	c.sendEvent(listCardsEvent)
	<-c.done
}

func (c *CardCLI) Stop() {
	c.FSM.SetState(string(endState))
}

func (c *CardCLI) cardMenu(e *fsm.Event) {
	if c.initTrack {
		if nextField, hasNext := c.activeField.Next(); hasNext {
			c.activeField = nextField
			c.sendEvent(editFieldEvent)

			return
		} else {
			c.initTrack = false
		}
	}

	c.ce.Print(NoField, -1)

	var (
		editTitleField       = "edit title"
		editDefinitionField  = "edit definition"
		editTranslationField = "edit translations"
		editExamplesField    = "edit examples"
		saveField            = "save"
		listField            = "list"
		removeField          = "delete"
		newField             = "new"
	)

	items := []string{
		editTitleField, editDefinitionField, editTranslationField, editExamplesField,
		saveField, removeField, listField, newField,
	}

	resultIdx, err := c.console.Select("Card", items)
	if err != nil {
		log.Error().Err(err).Msg("Prompt failed to get selection")
	}

	switch items[resultIdx] {
	case editTitleField:
		c.activeField = TitleField
		c.sendEvent(editFieldEvent)
	case editDefinitionField:
		c.activeField = DefinitionField
		c.sendEvent(editFieldEvent)
	case editTranslationField:
		c.activeField = TranslationField
		c.sendEvent(editFieldEvent)
	case editExamplesField:
		c.activeField = ExampleField
		c.sendEvent(selectExampleEvent)
	case saveField:
		c.sendEvent(saveCardEvent)
	case listField:
		c.sendEvent(listCardsEvent)
	case removeField:
		if err := c.ce.DeleteCard(); err != nil {
			c.console.Println("failed to delete card")
		}

		c.sendEvent(listCardsEvent)
	case newField:
		c.initTrack = true
		c.activeField = TitleField
		c.ce.ResetCard()
		c.FSM.SetState(string(editFieldState))
		c.sendEvent(writeFieldEvent)
	default:
		log.Error().Msgf("Invalid prompt index %d", resultIdx)
	}
}

func printMaxChars(s string, maxCount int) string {
	if len(s) <= maxCount {
		return s
	}

	return strings.TrimSpace(s[:maxCount]) + "..."
}

func (c *CardCLI) listCards(e *fsm.Event) {
	cardsList := c.ce.ListCards()

	items, err := formatCardsList(cardsList)
	if err != nil {
		log.Error().Err(err).Msg("Failed to format cards list")
	}

	selectedIndex, err := c.console.Select("Select card:", items)
	if err != nil {
		log.Warn().Err(err).Msg("unable to get selected card")
	}

	c.ce.SetCard(&cardsList[selectedIndex])
	c.sendEvent(editCardEvent)
}

func formatCardsList(cardsList []entity.Card) ([]string, error) {
	var tabs strings.Builder
	w := tabwriter.NewWriter(&tabs, 0, 0, 2, ' ', 0) //nolint:gomnd
	tFuncs := template.FuncMap{
		"printMax": printMaxChars,
		"join":     strings.Join,
	}

	for i, e := range cardsList {
		if i != 0 {
			fmt.Fprintf(w, "\n")
		}

		t := template.Must(template.New("cardsListItem").Funcs(tFuncs).Parse(cardListItem))

		err := t.Execute(w, e)
		if err != nil {
			return []string{}, err
		}
	}

	if err := w.Flush(); err != nil {
		return []string{}, err
	}

	return strings.Split(tabs.String(), "\n"), nil
}

func (c *CardCLI) showFieldPropositions(e *fsm.Event) {
	done := make(chan struct{})
	doneBusy := c.console.Busy(done)

	props, err := c.ce.GetPropositions(c.activeField, c.activeFieldIndex)

	close(done)

	<-doneBusy

	if err != nil {
		log.Warn().Err(err).Msg("unable to get propositions")
	}

	label := "Select " + c.activeField.String() + ":"

	index, err := c.console.Select(label, props)
	if err != nil {
		log.Warn().Err(err).Msg("unable to select proposition")
		c.sendEvent(editFieldEvent)

		return
	}

	err = c.ce.SetProposition(c.activeField, c.activeFieldIndex, index)
	if err != nil {
		log.Error().Err(err).Msg("failed to set proposition")
	}

	if c.activeField == ExampleField || c.activeField == TranslatedExampleField {
		c.sendEvent(editExampleEvent)

		return
	}

	c.sendEvent(editFieldEvent)
}

func (c *CardCLI) editFieldPrompt(e *fsm.Event) {
	c.ce.Print(c.activeField, c.activeFieldIndex)

	result, err := c.console.Prompt(c.activeField.String(), c.ce.GetField(c.activeField, c.activeFieldIndex))
	if err != nil {
		log.Warn().Err(err).Msg("error reading prompt")
		c.sendEvent(editFieldEvent)

		return
	}

	err = c.ce.SetField(c.activeField, c.activeFieldIndex, result)
	if err != nil {
		log.Warn().Err(err).Msg("error setting field")
		c.sendEvent(editFieldEvent)

		return
	}

	if c.activeField == ExampleField || c.activeField == TranslatedExampleField {
		c.sendEvent(editExampleEvent)

		return
	}

	c.sendEvent(editFieldEvent)
}

func (c *CardCLI) editFieldMenu(e *fsm.Event) {
	if c.activeField == CardField(0) {
		c.activeField = TitleField
	}

	c.ce.Print(c.activeField, 0)

	label := fmt.Sprintf("Edit field %s", c.activeField.String())

	var (
		showPropsField = "Show propositions"
		editField      = "Edit"
		validateField  = "Validate"
		cancelField    = "Cancel"
	)

	items := []string{editField, validateField, cancelField}

	if c.activeField != TitleField {
		items = append([]string{showPropsField}, items...)
	}

	resultIdx, err := c.console.Select(label, items)
	if err != nil {
		log.Error().Err(err).Msg("Prompt failed to get selection for editFieldMenu")

		if errors.Is(err, errInterrupt) || errors.Is(err, errEOF) {
			os.Exit(1)
		}
	}

	switch items[resultIdx] {
	case showPropsField:
		c.sendEvent(selectPropositionEvent)
	case editField:
		c.sendEvent(writeFieldEvent)
	case validateField:
		c.sendEvent(editCardEvent)
	case cancelField:
		c.initTrack = false
		c.sendEvent(editCardEvent)
	default:
		log.Error().Msgf("Invalid prompt index %d", resultIdx)
	}
}

func (c *CardCLI) editExampleMenu(e *fsm.Event) {
	c.ce.Print(c.activeField, c.activeFieldIndex)

	label := fmt.Sprintf("Edit field %s", c.activeField.String())

	var (
		showPropsField       = "Show sentences"
		editExampleField     = "Edit example"
		translateField       = "Translate"
		editTranslationField = "Edit translation"
		validateField        = "Validate"
		cancelField          = "Cancel"
		newField             = "New"
	)

	items := []string{showPropsField, editExampleField, translateField, editTranslationField, validateField, cancelField, newField}

	resultIdx, err := c.console.Select(label, items)
	if err != nil {
		log.Error().Err(err).Msg("Prompt failed to get selection for editExampleMenu")

		if errors.Is(err, errInterrupt) || errors.Is(err, errEOF) {
			os.Exit(1)
		}
	}

	switch items[resultIdx] {
	case showPropsField:
		// TODO
		c.activeField = ExampleField
		c.sendEvent(selectPropositionEvent)
	case editExampleField:
		c.activeField = ExampleField
		c.sendEvent(writeFieldEvent)
	case translateField:
		c.activeField = TranslatedExampleField
		c.sendEvent(selectPropositionEvent)
	case editTranslationField:
		c.activeField = TranslatedExampleField
		c.sendEvent(writeFieldEvent)
	case validateField:
		c.sendEvent(selectExampleEvent)
	case cancelField:
		c.sendEvent(selectExampleEvent)
	case newField:
		c.activeFieldIndex += 1
		c.FSM.SetState(string(selectExampleEvent))
		c.sendEvent(editFieldEvent)
	default:
		log.Error().Msgf("Invalid prompt index %d", resultIdx)
	}
}

func (c *CardCLI) selectExampleMenu(e *fsm.Event) {
	if len(c.ce.card.Examples) == 0 {
		c.sendEvent(editExampleEvent)
		c.activeFieldIndex = 0

		return
	}

	c.ce.Print(c.activeField, -1)

	label := "Select example"
	items := make([]string, 0, len(c.ce.card.Examples)+1)

	for i := range c.ce.card.Examples {
		items = append(items, strconv.Itoa(i+1))
	}

	items = append(items, "cancel")

	resultIdx, err := c.console.Select(label, items)
	if err != nil {
		log.Error().Err(err).Msg("Prompt failed to get selection for selectExampleMenu")

		if errors.Is(err, errInterrupt) || errors.Is(err, errEOF) {
			os.Exit(1)
		}
	}

	switch resultIdx {
	case len(items) - 1:
		// cancel
		c.sendEvent(editCardEvent)
	default:
		c.activeFieldIndex = resultIdx
		c.sendEvent(editExampleEvent)
	}
}

func (c *CardCLI) sendEvent(eventName event) {
	log.Trace().Msgf("can use event %s from state %s: %t. Available transitions: %s",
		eventName, c.FSM.Current(), c.FSM.Can(string(eventName)), c.FSM.AvailableTransitions())

	go func() {
		err := c.FSM.Event(string(eventName))
		if err != nil {
			log.Error().Err(err)
		}
	}()
}

func (c *CardCLI) save(e *fsm.Event) {
	if err := c.ce.SaveCard(); err != nil {
		log.Error().Err(err).Msg("failed to save card")
	}

	c.sendEvent(editCardEvent)
}
