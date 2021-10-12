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
	"github.com/thibaultmg/clingua/internal/usecase/language"
)

var getDefinitionsTimeout = 10 * time.Second
var getTranslationsTimeout = 15 * time.Second

type CardField int

const (
	NoField CardField = iota
	TitleField
	DefinitionField
	TranslationField
	ExempleField
	TranslatedExempleField
)

func (c CardField) String() string {
	switch c {
	case TitleField:
		return "Title"
	case DefinitionField:
		return "Definition"
	case TranslationField:
		return "Translations"
	case ExempleField:
		return "Exemples"
	default:
		// fmt.Println("card field", int(c))
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
		return ExempleField, true
	default:
		return NoField, false
	}
}

type CardEditor struct {
	card     *entity.Card
	cache    map[string]interface{}
	language language.LanguageUC
}

func NewCardEditor(card *entity.Card, lang language.LanguageUC) *CardEditor {
	return &CardEditor{
		card:     card,
		cache:    make(map[string]interface{}),
		language: lang,
	}
}

func (c *CardEditor) GetCard() *entity.Card {
	return c.card
}

func (c *CardEditor) GetField(field CardField, index int) string {
	switch field {
	case TitleField:
		return c.card.Title
	case DefinitionField:
		return c.card.Definition
	case TranslationField:
		if index >= len(c.card.Translations) {
			log.Debug().Msgf("invalid index %d for translations", index)
			return ""
		}
		return c.card.Translations[index]
	case ExempleField:
		if index >= len(c.card.Exemples) {
			log.Debug().Msgf("invalid index %d for exemples", index)
			return ""
		}
		return c.card.Exemples[index].Sentence
	case TranslatedExempleField:
		if index >= len(c.card.Exemples) {
			log.Debug().Msgf("invalid index %d for exemples", index)
			return ""
		}
		return c.card.Exemples[index].Translation
	default:
		return ""
	}
}

func (c *CardEditor) SetField(field CardField, index int, val string) error {
	switch field {
	case TitleField:
		c.card.Title = val
	case DefinitionField:
		c.card.Definition = val
	case TranslationField:
		if index > len(c.card.Translations) {
			log.Debug().Msgf("invalid index %d for translations", index)
		} else if index == len(c.card.Translations) || len(c.card.Translations) == 0 {
			c.card.Translations = append(c.card.Translations, val)
		} else {
			c.card.Translations[index] = val
		}
	case ExempleField:
		// TODO: handle index
		c.card.Exemples[index].Sentence = val
	case TranslatedExempleField:
		// TODO: handle index
		c.card.Exemples[index].Translation = val
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

	t.Execute(os.Stdout, c.card)
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
		cachekey := c.makeCacheKey(DefinitionField)
		cacheVal, ok := c.cache[cachekey]
		if !ok {
			panic("invalid cache key" + cachekey)
		}
		defProps, ok := cacheVal.([]language.DefinitionEntry)
		if !ok {
			panic("invalid type assertion")
		}
		c.SetField(DefinitionField, 0, defProps[index].Definition)
	case TranslationField:
		cachekey := c.makeCacheKey(TranslationField)
		cacheVal, ok := c.cache[cachekey]
		if !ok {
			panic("invalid cache key" + cachekey)
		}
		transProps, ok := cacheVal.([]string)
		if !ok {
			panic("invalid type assertion from card editor cache")
		}
		c.SetField(TranslationField, index, transProps[index])
	case ExempleField:
		return errors.New("not implemented")
	case TranslatedExempleField:
		return errors.New("not implemented")
	default:
		return errors.New("not implemented")
	}

	return nil
}

func (c *CardEditor) makeCacheKey(field CardField) string {
	return field.String() + "_" + c.GetField(field, 0)
}

func (c *CardEditor) getDefinitions() ([]string, error) {
	var defProps []language.DefinitionEntry
	var ret []string

	cachekey := c.makeCacheKey(DefinitionField)
	cacheVal, ok := c.cache[cachekey]
	if !ok {
		ctx, cancel := context.WithTimeout(context.Background(), getDefinitionsTimeout)
		defer cancel()

		var err error
		defProps, err = c.language.GetDefinition(ctx, c.card.Title, c.card.PartOfSpeech)
		if err != nil {
			return ret, fmt.Errorf("failed to get definitions: %v", err)
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
	err := t.Execute(&s, defProps)
	if err != nil {
		return ret, fmt.Errorf("unable to apply definition template: %v", err)
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
			return ret, fmt.Errorf("failed to get translations: %v", err)
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
