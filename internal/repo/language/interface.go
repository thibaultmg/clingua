package language

import (
	"context"

	"github.com/thibaultmg/clingua/internal/entity"
	"github.com/thibaultmg/clingua/internal/usecase/card"
	"golang.org/x/text/language"
)

type Dictionnary interface {
	GetDefinition(context.Context, string, language.Tag, entity.PartOfSpeech) ([]card.DefinitionEntry, error)
}

type Translator interface {
	GetTranslation(context.Context, string, entity.PartOfSpeech, string, string) ([]card.DefinitionEntry, error)
}
