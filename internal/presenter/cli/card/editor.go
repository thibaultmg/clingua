package card

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/thibaultmg/clingua/internal/entity"
	"github.com/thibaultmg/clingua/internal/usecase/language"
)

type CardField int

const (
	TitleField CardField = iota + 1
	DefinitionField
	TranslationsField
	ExemplesField
)

func (c CardField) getValue(card *entity.Card) string {
	switch c {
	case TitleField:
		return card.Title
	case DefinitionField:
		return card.Definition
	}

	return ""
}

func (c CardField) setValue(card *entity.Card, val string) {
	switch c {
	case TitleField:
		card.Title = val
	case DefinitionField:
		card.Definition = val
	}
}

func (c CardField) String() string {
	switch c {
	case TitleField:
		return "Title"
	case DefinitionField:
		return "Definition"
	case TranslationsField:
		return "Translations"
	case ExemplesField:
		return "Exemples"
	default:
		fmt.Println("card field", int(c))
		return ""
	}
}

type CardEditor struct {
	card     *entity.Card
	console  console
	language language.LanguageUC
	cache    map[string]interface{}
}

func NewCardEditor(card entity.Card, lang language.LanguageUC) *CardEditor {
	return &CardEditor{
		card:     &card,
		language: lang,
		console:  newConsole(os.Stdout),
		cache:    make(map[string]interface{}),
	}
}

func (c *CardEditor) GetCard() *entity.Card {
	return c.card
}

func (c *CardEditor) EditField(field CardField) error {

	result, err := c.console.Prompt(field.String(), field.getValue(c.card))
	if err != nil {
		return err
	}

	field.setValue(c.card, result)
	return nil
}

func (c *CardEditor) PrintField(field CardField) error {
	t := template.New(field.String())

	switch field {
	case DefinitionField:
		t = template.Must(t.Parse(definitionTemplate))
	case TitleField:
		t = template.Must(t.Parse(titleTemplate))
	}

	return t.Execute(os.Stdout, c.card)
}

func (c *CardEditor) SelectProposition(field CardField) error {
	// fmt.Println("coucou from showFieldPropositions")
	switch field {
	case DefinitionField:
		val, err := c.proposeDefinitions()
		if err != nil {
			log.Warn().Msgf("Unable to fetch definitions: %v", err)
			// c.sendEvent(setFieldEvent.String())
			return err
		}
		c.card.Definition = val
	}

	return nil
}

func (c *CardEditor) proposeDefinitions() (string, error) {
	var def []language.DefinitionEntry

	// // To Delete
	// mockDef := []language.DefinitionEntry{
	// 	{
	// 		PartOfSpeech: entity.Interjection,
	// 		Definition:   "this is a definition",
	// 		Registers:    []string{"slang", "formal"},
	// 		Domains:      []string{"sport", "tennis"},
	// 		Provider:     "oxford",
	// 	},
	// 	{
	// 		PartOfSpeech: entity.Any,
	// 		Definition:   "this is a definition",
	// 	},
	// }
	// c.cache["definition"] = mockDef

	// Check cache
	cacheVal, ok := c.cache["definition"]
	if !ok {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		c.console.Busy(ctx.Done())
		var err error
		def, err = c.language.GetDefinition(ctx, c.card.Title, c.card.PartOfSpeech)
		if err != nil {
			return "", fmt.Errorf("failed to get definitions: %v", err)
		}
		// Close channel for busy
		cancel()

		// Set cache
		c.cache["definition"] = def
	} else {
		def, ok = cacheVal.([]language.DefinitionEntry)
		if !ok {
			panic("invalid type assertion from card editor cache")
		}
	}

	funcs := template.FuncMap{"join": strings.Join}
	funcs["add"] = func(i, j int) int { return i + j }
	t := template.Must(template.New("definitions").Funcs(funcs).Parse(definitionPropsTemplate))
	var s strings.Builder
	t.Execute(&s, def)
	splittedDefs := strings.Split(s.String(), "\n")
	items := make([]string, 0, len(splittedDefs))
	for _, val := range splittedDefs {
		if len(strings.TrimSpace(val)) > 0 {
			items = append(items, val)
		}
	}

	label := "Select the definition"
	if len(def) > 0 && len(def[0].Provider) > 0 {
		label = label + " from " + def[0].Provider
	}
	resultIdx, err := c.console.Select(label, items)
	if err != nil {
		return "", fmt.Errorf("failed to get selected definition from prompt: %v", err)
	}

	if len(def) <= resultIdx {
		return "", fmt.Errorf("invalid index in slice of definitions propositions: %d", resultIdx)
	}

	return def[resultIdx].Definition, nil
}
