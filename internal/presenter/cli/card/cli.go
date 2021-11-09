package card

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/looplab/fsm"
	"github.com/rs/zerolog/log"
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
)

type CardCLI struct {
	ce          *CardEditor
	FSM         *fsm.FSM
	console     console
	activeField CardField
	initTrack   bool
	done        <-chan struct{}
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
			{Name: string(editCardEvent), Src: []string{string(listCardsState), string(editFieldState)}, Dst: string(editCardState)},
			{Name: string(selectPropositionEvent), Src: []string{string(editFieldState)}, Dst: string(selectPropositionState)},
			{Name: string(writeFieldEvent), Src: []string{string(editFieldState)}, Dst: string(writeFieldState)},
			{Name: string(quitEvent), Src: []string{string(editCardState), string(saveCardState)}, Dst: string(endState)},
			{Name: string(listCardsEvent), Src: []string{string(editCardState), string(initState)}, Dst: string(listCardsState)},
			{Name: string(saveCardEvent), Src: []string{string(editCardState)}, Dst: string(saveCardState)},
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
	c.sendEvent(string(editFieldState))
	<-c.done
}

func (c *CardCLI) RunList() {
	c.sendEvent(string(listCardsEvent))
	<-c.done
}

func (c *CardCLI) Stop() {
	c.FSM.SetState(string(endState))
}

func (c *CardCLI) cardMenu(e *fsm.Event) {
	if c.initTrack {
		if nextField, hasNext := c.activeField.Next(); hasNext {
			c.activeField = nextField
			c.sendEvent(string(editFieldEvent))

			return
		} else {
			c.initTrack = false
		}
	}

	c.ce.Print(NoField)

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
		c.sendEvent(string(editFieldEvent))
	case editDefinitionField:
		c.activeField = DefinitionField
		c.sendEvent(string(editFieldEvent))
	case editTranslationField:
		c.activeField = TranslationField
		c.sendEvent(string(editFieldEvent))
	case editExamplesField:
		c.activeField = ExampleField
		c.sendEvent(string(editFieldEvent))
	case saveField:
		c.sendEvent(string(saveCardEvent))
	case listField:
		c.sendEvent(string(listCardsEvent))
	case removeField:
		if err := c.ce.DeleteCard(); err != nil {
			c.console.Println("failed to delete card")
		}

		c.sendEvent(string(listCardsEvent))
	case newField:
		c.initTrack = true
		c.activeField = TitleField
		c.ce.ResetCard()
		c.FSM.SetState(string(initState))
		c.sendEvent(string(editFieldState))
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
	// items := make([]string, 0, len(cardsList))
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
			log.Error().Err(err).Msg("failed to execute template on cards list item")
		}
	}

	if err := w.Flush(); err != nil {
		log.Warn().Err(err).Msg("Failed to flush tabwriter")
	}

	items := strings.Split(tabs.String(), "\n")

	selectedIndex, err := c.console.Select("Select card:", items)
	if err != nil {
		log.Warn().Err(err).Msg("unable to get selected card")
	}

	c.ce.SetCard(&cardsList[selectedIndex])
	c.sendEvent(string(editCardEvent))
}

func (c *CardCLI) showFieldPropositions(e *fsm.Event) {
	done := make(chan struct{})
	doneBusy := c.console.Busy(done)

	props, err := c.ce.GetPropositions(c.activeField, 0)

	close(done)

	<-doneBusy

	if err != nil {
		log.Warn().Err(err).Msg("unable to get propositions")
	}

	label := "Select " + c.activeField.String() + ":"

	index, err := c.console.Select(label, props)
	if err != nil {
		log.Warn().Err(err).Msg("unable to select proposition")
		c.sendEvent(string(editFieldEvent))

		return
	}

	err = c.ce.SetProposition(c.activeField, index)
	if err != nil {
		log.Error().Err(err).Msg("failed to set proposition")
	}

	c.sendEvent(string(editFieldEvent))
}

func (c *CardCLI) editFieldPrompt(e *fsm.Event) {
	c.ce.Print(c.activeField)

	result, err := c.console.Prompt(c.activeField.String(), c.ce.GetField(c.activeField, 0))
	if err != nil {
		log.Warn().Err(err).Msg("error reading prompt")
		c.sendEvent(string(editFieldEvent))

		return
	}

	err = c.ce.SetField(c.activeField, 0, result)
	if err != nil {
		log.Warn().Err(err).Msg("error setting field")
		c.sendEvent(string(editFieldEvent))

		return
	}

	c.sendEvent(string(editFieldEvent))
}

func (c *CardCLI) editFieldMenu(e *fsm.Event) {
	if c.activeField == CardField(0) {
		c.activeField = TitleField
	}

	c.ce.Print(c.activeField)

	label := fmt.Sprintf("Edit field %s", c.activeField.String())

	var (
		showPropsField = "Show propositions"
		editField      = "Edit"
		validateField  = "Validate"
		cancelField    = "Cancel"
	)

	items := []string{showPropsField, editField, validateField, cancelField}

	resultIdx, err := c.console.Select(label, items)
	if err != nil {
		log.Error().Err(err).Msg("Prompt failed to get selection for editFieldMenu")

		if errors.Is(err, errInterrupt) || errors.Is(err, errEOF) {
			os.Exit(1)
		}
	}

	switch items[resultIdx] {
	case showPropsField:
		c.sendEvent(string(selectPropositionEvent))
	case editField:
		c.sendEvent(string(writeFieldEvent))
	case validateField:
		c.sendEvent(string(editCardEvent))
	case cancelField:
		c.initTrack = false
		c.sendEvent(string(editCardEvent))
	default:
		log.Error().Msgf("Invalid prompt index %d", resultIdx)
	}
}

func (c *CardCLI) sendEvent(eventName string) {
	log.Trace().Msgf("can use event %s from state %s: %t. Available transitions: %s",
		eventName, c.FSM.Current(), c.FSM.Can(eventName), c.FSM.AvailableTransitions())

	go func() {
		err := c.FSM.Event(eventName)
		if err != nil {
			log.Error().Err(err)
		}
	}()
}

func (c *CardCLI) save(e *fsm.Event) {
	if err := c.ce.SaveCard(); err != nil {
		log.Error().Err(err).Msg("failed to save card")
	}

	c.sendEvent(string(editCardEvent))
}
