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
	dict   Dictionary
	trans  Translator
	wtrans WordTranslator
	sent   Sentences
}

func New(dict Dictionary, trans Translator, wtrans WordTranslator) *LanguageUCImpl {
	return &LanguageUCImpl{
		dict:   dict,
		trans:  trans,
		wtrans: wtrans,
	}
}

func (l *LanguageUCImpl) Define(ctx context.Context, word string, pos entity.PartOfSpeech) ([]DefinitionEntry, error) {
	log.Debug().Msgf("Getting definition for word %s with pos %v", word, pos)

	resp, err := l.dict.Define(ctx, word, pos)
	if err != nil {
		log.Error().Err(err).Str("word", word).Str("partOfSpeech", pos.String()).Msg("Unable to get definition")

		return resp, err
	}

	return resp, nil
}

func (l *LanguageUCImpl) Translate(ctx context.Context, text string) ([]string, error) {
	log.Debug().Msgf("Getting translation for text %q", text)

	resp, err := l.trans.Translate(ctx, text)
	if err != nil {
		log.Error().Err(err).Str("text", text).Msg("Unable to get translations")

		return resp, err
	}

	return resp, nil
}

func (l *LanguageUCImpl) TranslateWord(ctx context.Context, word string, pos entity.PartOfSpeech) ([]WordTranslationEntry, error) {
	log.Debug().Msgf("Getting word translation for %q with pos %q", word, pos)

	resp, err := l.wtrans.TranslateWord(ctx, word, pos)
	if err != nil {
		log.Error().Err(err).Str("word", word).Str("partOfSpeech", pos.String()).Msg("Unable to get word translations")

		return resp, err
	}

	return resp, nil
}

func (l *LanguageUCImpl) Sentences(ctx context.Context, text string) ([]string, error) {
	log.Debug().Msgf("Getting sentences examples for text %q", text)

	resp, err := l.sent.Sentences(ctx, text)
	if err != nil {
		log.Error().Err(err).Str("text", text).Msg("Unable to get sentences")

		return resp, err
	}

	return resp, nil
}
