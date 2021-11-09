package filesystem

import (
	"github.com/rs/zerolog/log"
	"golang.org/x/text/language"

	"github.com/thibaultmg/clingua/internal/entity"
)

type card struct {
	ID                   string
	FromLanguage         string
	ToLanguage           string
	Title                string
	PartOfSpeech         string
	Register             string
	Definition           string
	Translations         []string
	Examples             []string
	ExamplesTranslations []string
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

	return entity.Card{
		ID:                   c.ID,
		From:                 from,
		To:                   to,
		Title:                c.Title,
		PartOfSpeech:         pos,
		Register:             c.Register,
		Definition:           c.Definition,
		Examples:             c.Examples,
		ExamplesTranslations: c.ExamplesTranslations,
		Translations:         c.Translations,
	}
}

func entityToCard(ecard *entity.Card) card {
	return card{
		ID:                   ecard.ID,
		FromLanguage:         ecard.From.String(),
		ToLanguage:           ecard.To.String(),
		Title:                ecard.Title,
		PartOfSpeech:         ecard.PartOfSpeech.String(),
		Register:             ecard.Register,
		Definition:           ecard.Definition,
		Translations:         ecard.Translations,
		Examples:             ecard.Examples,
		ExamplesTranslations: ecard.ExamplesTranslations,
	}
}
