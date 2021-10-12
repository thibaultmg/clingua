package language

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/thibaultmg/clingua/internal/entity"
)

// Language struct acts as a facade integrating all language external services into a simplified interface.
// A language instance is dedicated to one fromLanguage (the learner native language) and one toLanguage
// (the learned language).
type LanguageUCImpl struct {
	dict  Dictionnary
	trans Translator
}

func New(dict Dictionnary, trans Translator) *LanguageUCImpl {
	return &LanguageUCImpl{
		dict:  dict,
		trans: trans,
	}
}

func (l *LanguageUCImpl) GetDefinition(ctx context.Context, word string, pos entity.PartOfSpeech) ([]DefinitionEntry, error) {
	log.Debug().Msgf("Getting definition for word %s with pos %v", word, pos)

	resp, err := l.dict.GetDefinition(ctx, word, pos)
	if err != nil {
		log.Error().Err(err).Str("word", word).Str("partOfSpeech", pos.String()).Msg("Unable to get definition")
		return resp, err
	}

	return resp, nil
}

func (l *LanguageUCImpl) GetTranslation(ctx context.Context, word string, pos entity.PartOfSpeech) ([]string, error) {
	log.Debug().Msgf("Getting translation for text %s with pos %v", word, pos)

	resp, err := l.trans.GetTranslation(ctx, word, pos)
	if err != nil {
		log.Error().Err(err).Str("word", word).Str("partOfSpeech", pos.String()).Msg("Unable to get definition")
		return resp, err
	}

	return resp, nil
}
