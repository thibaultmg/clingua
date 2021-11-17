package filesystem

import (
	"github.com/rs/zerolog/log"
	"golang.org/x/text/language"

	"github.com/thibaultmg/clingua/internal/entity"
)

type card struct {
	ID           string
	FromLanguage string
	ToLanguage   string
	Title        string
	PartOfSpeech string
	Register     string
	Definition   string
	Translations []string
	Examples     []example
}

type example struct {
	Example     string
	Translation string
}

func (c card) ToEntity() entity.Card {
	var (
		pos entity.PartOfSpeech
		err error
	)

	if len(c.PartOfSpeech) > 0 {
		pos, err = entity.ParsePartOfSpeech(c.PartOfSpeech)
		if err != nil {
			log.Err(err).Msg("invalid part of speech")
		}
	}

	from := language.MustParse(c.FromLanguage)
	to := language.MustParse(c.ToLanguage)

	examples := make([]entity.Example, 0, len(c.Examples))
	for _, e := range c.Examples {
		examples = append(examples, entity.Example{Example: e.Example, Translation: e.Translation})
	}

	return entity.Card{
		ID:           c.ID,
		From:         from,
		To:           to,
		Title:        c.Title,
		PartOfSpeech: pos,
		Register:     c.Register,
		Definition:   c.Definition,
		Examples:     examples,
		Translations: c.Translations,
	}
}

func entityToCard(ecard *entity.Card) card {
	examples := make([]example, 0, len(ecard.Examples))
	for _, e := range ecard.Examples {
		examples = append(examples, example{Example: e.Example, Translation: e.Translation})
	}

	return card{
		ID:           ecard.ID,
		FromLanguage: ecard.From.String(),
		ToLanguage:   ecard.To.String(),
		Title:        ecard.Title,
		PartOfSpeech: ecard.PartOfSpeech.String(),
		Register:     ecard.Register,
		Definition:   ecard.Definition,
		Translations: ecard.Translations,
		Examples:     examples,
	}
}
