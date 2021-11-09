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
	getDefinitionsTimeout  = 10 * time.Second
	getTranslationsTimeout = 15 * time.Second
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
	cache    map[string]interface{}
	language language.LanguageUC
	cardUC   card.CardUC
}

func NewCardEditor(card *entity.Card, lang language.LanguageUC, cardUC card.CardUC) *CardEditor {
	return &CardEditor{
		card:     card,
		cache:    make(map[string]interface{}),
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
		cacheVal := c.mustGetCacheVal(DefinitionField)

		defProps, ok := cacheVal.([]language.DefinitionEntry)
		if !ok {
			panic("invalid type assertion")
		}

		err := c.SetField(DefinitionField, 0, defProps[index].Definition)
		if err != nil {
			log.Error().Err(err).Msg("failed to set field on card")
		}
	case TranslationField:
		cacheVal := c.mustGetCacheVal(TranslationField)

		transProps, ok := cacheVal.([]string)
		if !ok {
			panic("invalid type assertion from card editor cache")
		}

		err := c.SetField(TranslationField, index, transProps[index])
		if err != nil {
			log.Error().Err(err).Msg("failed to set field on card")
		}
	case ExampleField, TranslatedExampleField, NoField, TitleField:
		return errors.New("not implemented")
	default:
		return errors.New("not implemented")
	}

	return nil
}

func (c *CardEditor) mustGetCacheVal(field CardField) interface{} {
	cachekey := c.makeCacheKey(field)

	cacheVal, ok := c.cache[cachekey]
	if !ok {
		panic("invalid cache key" + cachekey)
	}

	return cacheVal
}

func (c *CardEditor) makeCacheKey(field CardField) string {
	switch field {
	case DefinitionField, TranslationField, ExampleField:
		return field.String() + "_" + c.GetField(TitleField, 0)
	default:
		log.Error().Msg("cache key not implemented")

		return field.String() + "_" + c.GetField(field, 0)
	}
}

func (c *CardEditor) getCachedDefinitions() ([]language.DefinitionEntry, error) {
	var ret []language.DefinitionEntry

	cachekey := c.makeCacheKey(DefinitionField)

	cacheVal, ok := c.cache[cachekey]
	if !ok {
		return ret, errors.New("failed to get cached definitions, no cache key")
	}

	ret, ok = cacheVal.([]language.DefinitionEntry)
	if !ok {
		return ret, errors.New("failed to cast cached definitions")
	}

	return ret, nil
}

func (c *CardEditor) getDefinitions() ([]string, error) {
	var (
		defProps []language.DefinitionEntry
		ret      []string
	)

	cachekey := c.makeCacheKey(DefinitionField)

	cacheVal, ok := c.cache[cachekey]
	if !ok {
		ctx, cancel := context.WithTimeout(context.Background(), getDefinitionsTimeout)
		defer cancel()

		var err error

		defProps, err = c.language.GetDefinition(ctx, c.card.Title, c.card.PartOfSpeech)
		if err != nil {
			return ret, fmt.Errorf("failed to get definitions: %w", err)
		}

		// Set cache
		c.cache[cachekey] = defProps
	} else {
		defProps, ok = cacheVal.([]language.DefinitionEntry)
		if !ok {
			panic("invalid type assertion")
		}
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
	var ret []string

	cachekey := c.makeCacheKey(TranslationField)

	cacheVal, ok := c.cache[cachekey]
	if !ok {
		ctx, cancel := context.WithTimeout(context.Background(), getTranslationsTimeout)
		defer cancel()

		var err error

		ret, err = c.language.GetTranslation(ctx, c.card.Title, c.card.PartOfSpeech)
		if err != nil {
			return ret, fmt.Errorf("failed to get translations: %w", err)
		}

		// Set cache
		c.cache[cachekey] = ret
	} else {
		ret, ok = cacheVal.([]string)
		if !ok {
			panic("invalid type assertion from card editor cache")
		}
	}

	return ret, nil
}

func (c *CardEditor) setDefinitionField(val string, index int) {
	c.card.Definition = val

	defs, err := c.getCachedDefinitions()
	if err != nil {
		log.Error().Err(err).Msg("failed to get cached definition")
	} else {
		c.card.PartOfSpeech = defs[index].PartOfSpeech
	}
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
