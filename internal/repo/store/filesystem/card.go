package filesystem

import (
	"github.com/rs/zerolog/log"
	"github.com/thibaultmg/clingua/internal/entity"
)

type card struct {
	FromLanguage         string
	ToLanguage           string
	Title                string
	PartOfSpeech         string
	Register             string
	Definition           string
	Translations         []string
	Exemples             []string
	ExemplesTranslations []string
}

func (c card) ToEntity() entity.Card {
	ret := entity.NewCard()
	ret.Definition = c.Definition
	ret.Title = c.Title
	if len(c.PartOfSpeech) > 0 {
		pos, err := entity.ParsePartOfSpeech(c.PartOfSpeech)
		if err != nil {
			log.Err(err).Msg("invalid part of speech")
		}

		ret.PartOfSpeech = pos
	}

	return ret
}
