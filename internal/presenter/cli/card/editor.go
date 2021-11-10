package card

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/thibaultmg/clingua/internal/entity"
	"github.com/thibaultmg/clingua/internal/usecase/card"
	"github.com/thibaultmg/clingua/internal/usecase/language"
)

var (
	getDefinitionsTimeout   = 10 * time.Second
	getTranslationsTimeout  = 30 * time.Second
	sentenceExampleMinWords = 3
)

type CardField int

const (
	NoField CardField = iota
	TitleField
	DefinitionField
	TranslationField
	ExampleField
	TranslatedExampleField
)

func (c CardField) String() string {
	switch c {
	case TitleField:
		return "Title"
	case DefinitionField:
		return "Definition"
	case TranslationField:
		return "Translations"
	case ExampleField:
		return "Examples"
	case TranslatedExampleField:
		return "Example translation"
	default:
		return ""
	}
}

func (c CardField) Next() (CardField, bool) {
	switch c {
	case TitleField:
		return DefinitionField, true
	case DefinitionField:
		return TranslationField, true
	case TranslationField:
		return ExampleField, true
	default:
		return NoField, false
	}
}

type CardEditor struct {
	card     *entity.Card
	language language.LanguageUC
	cardUC   card.CardUC
}

func NewCardEditor(card *entity.Card, lang language.LanguageUC, cardUC card.CardUC) *CardEditor {
	return &CardEditor{
		card:     card,
		language: lang,
		cardUC:   cardUC,
	}
}

func (c *CardEditor) GetCard() *entity.Card {
	return c.card
}

func (c *CardEditor) SetCard(card *entity.Card) {
	c.card = card
}

func (c *CardEditor) ResetCard() {
	newCard := entity.NewCard()
	c.card = &newCard
}

func (c *CardEditor) SaveCard() error {
	_, err := c.cardUC.Create(context.Background(), *c.card)

	return err
}

func (c *CardEditor) DeleteCard() error {
	return c.cardUC.Delete(context.Background(), c.card.ID)
}

func (c *CardEditor) ListCards() []entity.Card {
	cardsList, err := c.cardUC.List(context.Background())
	if err != nil {
		log.Error().Err(err)
	}

	return cardsList
}

func (c *CardEditor) GetField(field CardField, index int) string {
	switch field {
	case TitleField:
		return c.card.Title
	case DefinitionField:
		return c.card.Definition
	case TranslationField, ExampleField, TranslatedExampleField:
		return c.getFieldWithIndex(field, index)
	default:
		log.Warn().Msgf("invalid get field %v", field)

		return ""
	}
}

func (c *CardEditor) SetField(field CardField, index int, val string) error {
	switch field {
	case TitleField:
		c.card.Title = val
	case DefinitionField:
		c.setDefinitionField(val, index)
	case TranslationField:
		c.setTranslationField(val, index)
	case ExampleField:
		// TODO: handle index
		c.card.Examples[index] = val
	case TranslatedExampleField:
		// TODO: handle index
		c.card.ExamplesTranslations[index] = val
	default:
		return fmt.Errorf("invalid index %d on field %s", index, field)
	}

	return nil
}

func (c CardEditor) Print(field CardField) {
	tFuncs := template.FuncMap{
		"join": strings.Join,
	}

	var t *template.Template

	switch field {
	case NoField:
		t = template.Must(template.New("card").Funcs(tFuncs).Parse(cardTemplate))
	case TitleField:
		t = template.Must(template.New("title").Parse(titleTemplate))
	case DefinitionField:
		t = template.Must(template.New("definition").Parse(definitionTemplate))
	case TranslationField:
		t = template.Must(template.New("translation").Funcs(tFuncs).Parse(translationTemplate))
	}

	err := t.Execute(os.Stdout, c.card)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute template for printing card field")
	}
}

func (c *CardEditor) GetPropositions(field CardField, index int) ([]string, error) {
	switch field {
	case DefinitionField:
		return c.getDefinitions()
	case TranslationField:
		return c.getTranslations()
	default:
		return []string{}, fmt.Errorf("invalid field %s or index %d", field, index)
	}
}

func (c *CardEditor) SetProposition(field CardField, index int) error {
	switch field {
	case DefinitionField:
		defProps, err := c.language.Define(context.Background(), c.card.Title, c.card.PartOfSpeech)
		if err != nil {
			return fmt.Errorf("failed to get definitions: %w", err)
		}

		c.card.Definition = defProps[index].Definition
		c.card.PartOfSpeech = defProps[index].PartOfSpeech

		for _, e := range defProps[index].Examples {
			if len(strings.Fields(e)) >= sentenceExampleMinWords {
				c.card.Examples = append(c.card.Examples, e)
			}
		}
	case TranslationField:
		transRes, err := c.language.TranslateWord(context.Background(), c.card.Title, c.card.PartOfSpeech)
		if err != nil {
			return fmt.Errorf("failed to get translations: %w", err)
		}

		c.card.Translations = append(c.card.Translations, transRes[index].Translation)
	case ExampleField, TranslatedExampleField, NoField, TitleField:
		return errors.New("not implemented")
	default:
		return errors.New("not implemented")
	}

	return nil
}

func (c *CardEditor) getDefinitions() ([]string, error) {
	var (
		defProps []language.DefinitionEntry
		ret      []string
	)

	ctx, cancel := context.WithTimeout(context.Background(), getDefinitionsTimeout)
	defer cancel()

	var err error

	defProps, err = c.language.Define(ctx, c.card.Title, c.card.PartOfSpeech)
	if err != nil {
		return ret, fmt.Errorf("failed to get definitions: %w", err)
	}

	// Format props for screen print
	var s strings.Builder

	tFuncs := template.FuncMap{
		"join": strings.Join,
		"add":  func(i, j int) int { return i + j },
	}
	t := template.Must(template.New("definitions").Funcs(tFuncs).Parse(definitionPropsTemplate))

	if err := t.Execute(&s, defProps); err != nil {
		return ret, fmt.Errorf("unable to apply definition template: %w", err)
	}

	ret = make([]string, 0, len(defProps))

	splittedDefs := strings.Split(s.String(), "\n")
	for _, val := range splittedDefs {
		if len(strings.TrimSpace(val)) > 0 {
			ret = append(ret, val)
		}
	}

	return ret, nil
}

func (c *CardEditor) getTranslations() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), getTranslationsTimeout)
	defer cancel()

	transRes, err := c.language.TranslateWord(ctx, c.card.Title, c.card.PartOfSpeech)
	if err != nil {
		return []string{}, fmt.Errorf("failed to get translations: %w", err)
	}

	ret := make([]string, 0, len(transRes))

	for _, e := range transRes {
		ret = append(ret, fmt.Sprintf("%s\t— %s [%s]", e.Translation, e.PartOfSpeech, e.Meaning))
	}

	return ret, nil
}

func (c *CardEditor) setDefinitionField(val string, index int) {
	c.card.Definition = val

	defProps, err := c.language.Define(context.Background(), c.card.Title, c.card.PartOfSpeech)
	if err != nil {
		log.Error().Err(err).Msg("failed to get definitions to set value in card")
	}

	c.card.PartOfSpeech = defProps[index].PartOfSpeech
}

func (c *CardEditor) setTranslationField(val string, index int) {
	switch transCount := len(c.card.Translations); {
	case index > transCount:
		log.Debug().Msgf("invalid index %d for translations", index)
	case transCount == index, transCount == 0:
		c.card.Translations = append(c.card.Translations, val)
	default:
		c.card.Translations[index] = val
	}
}

func (c *CardEditor) getFieldWithIndex(fieldType CardField, index int) string {
	var fieldValue []string

	switch fieldType {
	case TranslationField:
		fieldValue = c.card.Translations
	case ExampleField:
		fieldValue = c.card.Examples
	case TranslatedExampleField:
		fieldValue = c.card.ExamplesTranslations
	default:
		log.Error().Msgf("invalid indexed field type %s", fieldType)

		return ""
	}

	if index >= len(fieldValue) {
		log.Error().Msgf("invalid index %d for field %v", index, fieldType)

		return ""
	}

	return fieldValue[index]
}
