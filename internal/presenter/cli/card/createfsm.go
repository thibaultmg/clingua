package card

import (
	"errors"
	"fmt"
	"os"
	"text/template"

	"github.com/looplab/fsm"
	"github.com/rs/zerolog/log"
)

type event string

func (e event) String() string {
	return string(e)
}

type state string

func (s state) String() string {
	return string(s)
}

const (
	startCreationState     state = "startCreation"
	editFieldMenuState     state = "editFieldMenu"
	cardMenuState          state = "cardMenu"
	saveState              state = "save"
	endState               state = "end"
	fieldPropositionsState state = "fieldPropositions"
	editFieldPromptState   state = "editFieldPrompt"
)

const (
	startEvent                 event = "start"
	cancelEditFieldEvent       event = "cancelEditField"
	quitEvent                  event = "quit"
	editFieldEvent             event = "editField"
	showFieldPropositionsEvent event = "showFieldPropositions"
	setFieldEvent              event = "setField"
	validateFieldEvent         event = "validateField"
	saveCardEvent              event = "saveCard"
	showCardEvent              event = "showCard"
	nextFieldEvent             event = "nextField"
)

type CardFSM struct {
	editor           *CardEditor
	FSM              *fsm.FSM
	console          console
	activeField      CardField
	activeFieldIndex int
	done             <-chan struct{}
}

func NewCardFSM(editor *CardEditor) *CardFSM {
	doneChan := make(chan struct{})
	ret := &CardFSM{
		editor:  editor,
		done:    doneChan,
		console: newConsole(os.Stdout),
	}

	ret.FSM = fsm.NewFSM(
		"",
		fsm.Events{
			{Name: startEvent.String(), Src: []string{startCreationState.String()}, Dst: editFieldMenuState.String()},
			{Name: cancelEditFieldEvent.String(), Src: []string{editFieldMenuState.String()}, Dst: cardMenuState.String()},
			{Name: quitEvent.String(), Src: []string{cardMenuState.String(), saveState.String()}, Dst: endState.String()},
			{Name: editFieldEvent.String(), Src: []string{editFieldMenuState.String()}, Dst: editFieldPromptState.String()},
			{Name: showFieldPropositionsEvent.String(), Src: []string{editFieldMenuState.String()}, Dst: fieldPropositionsState.String()},
			{Name: setFieldEvent.String(), Src: []string{editFieldPromptState.String(), fieldPropositionsState.String(), cardMenuState.String()}, Dst: editFieldMenuState.String()},
			{Name: validateFieldEvent.String(), Src: []string{editFieldMenuState.String()}, Dst: cardMenuState.String()},
			{Name: saveCardEvent.String(), Src: []string{cardMenuState.String()}, Dst: saveState.String()},
			{Name: showCardEvent.String(), Src: []string{cardMenuState.String()}, Dst: cardMenuState.String()},
			{Name: nextFieldEvent.String(), Src: []string{editFieldMenuState.String()}, Dst: editFieldMenuState.String()},
		},
		fsm.Callbacks{
			"enter_state": func(e *fsm.Event) {
				log.Debug().Msgf("Enter state with event %s from state %s to state %s", e.Event, e.Src, e.Dst)
			},
			"enter_" + editFieldMenuState.String():     ret.editFieldMenu,
			"enter_" + fieldPropositionsState.String(): ret.showFieldPropositions,
			"enter_" + editFieldPromptState.String():   ret.editFieldPrompt,
			"enter_" + cardMenuState.String():          ret.cardMenu,
			"enter_" + saveState.String():              ret.save,
			"enter_" + endState.String(): func(e *fsm.Event) {
				close(doneChan)
			},
		},
	)

	log.Debug().Msg(fsm.Visualize(ret.FSM))

	return ret
}

func (c *CardFSM) Run() {
	c.FSM.SetState(startCreationState.String())
	c.activeField = DefinitionField
	go c.sendEvent(startEvent.String())
	<-c.done
}

func (c *CardFSM) Stop() {
	c.FSM.SetState(endState.String())
}

func (c *CardFSM) cardMenu(e *fsm.Event) {
	items := []string{"edit title", "edit definition", "edit translations", "edit exemples", "validate", "cancel"}
	resultIdx, err := c.console.Select("Card", items)
	if err != nil {
		log.Error().Err(err).Msg("Prompt failed to get selection")
	}

	switch resultIdx {
	case 0:
		c.activeField = TitleField
		go c.sendEvent(editFieldEvent.String())
	case 1:
		c.activeField = DefinitionField
		go c.sendEvent(editFieldEvent.String())
	case 2:
		c.activeField = TranslationsField
		go c.sendEvent(editFieldEvent.String())
	case 3:
		c.activeField = ExemplesField
		go c.sendEvent(editFieldEvent.String())
	case 4:
		go c.sendEvent(saveCardEvent.String())
	case 5:
		go c.sendEvent(quitEvent.String())
	default:
		log.Error().Msgf("Invalid prompt index %d", resultIdx)
	}
}

func (c *CardFSM) showFieldPropositions(e *fsm.Event) {
	err := c.editor.SelectProposition(c.activeField)
	if err != nil {
		log.Warn().Err(err).Msg("unable to use propositions")
	}

	go c.sendEvent(setFieldEvent.String())
}

func (c *CardFSM) editFieldPrompt(e *fsm.Event) {
	err := c.editor.EditField(c.activeField)
	if err != nil {
		log.Warn().Err(err).Msg("unable to edit field")
	}

	go c.sendEvent(setFieldEvent.String())
}

func (c *CardFSM) editFieldMenu(e *fsm.Event) {
	if c.activeField == CardField(0) {
		c.activeField = TitleField
	}

	t := template.Must(template.New("definition").Parse(definitionTemplate))
	t.Execute(os.Stdout, c.editor.GetCard())

	label := fmt.Sprintf("Edit field %s", c.activeField.String())
	items := []string{"Show propositions", "Edit", "Validate", "Cancel"}
	resultIdx, err := c.console.Select(label, items)
	if err != nil {
		log.Error().Err(err).Msg("Prompt failed to get selection for editFieldMenu")
		if errors.Is(err, errInterrupt) || errors.Is(err, errEOF) {
			os.Exit(1)
		}
	}

	switch resultIdx {
	case 0:
		go c.sendEvent(showFieldPropositionsEvent.String())
	case 1:
		go c.sendEvent(editFieldEvent.String())
	case 2:
		go c.sendEvent(validateFieldEvent.String())
	case 3:
		go c.sendEvent(cancelEditFieldEvent.String())
	default:
		log.Error().Msgf("Invalid prompt index %d", resultIdx)
	}
}

func (c *CardFSM) sendEvent(eventName string) {
	log.Debug().Msgf("can use event %s from state %s: %t. Available transitions: %#s", eventName, c.FSM.Current(), c.FSM.Can(eventName), c.FSM.AvailableTransitions())

	err := c.FSM.Event(eventName)
	if err != nil {
		log.Error().Err(err)
	}
}

func (c *CardFSM) save(e *fsm.Event) {
	fmt.Println("Saving card")
	go c.sendEvent(quitEvent.String())
	os.Exit(0)
}
