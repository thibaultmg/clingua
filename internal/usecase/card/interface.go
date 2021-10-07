package card

import (
	"context"

	"github.com/thibaultmg/clingua/internal/entity"
)

type DefinitionRepo interface {
	Get(context.Context, string, entity.PartOfSpeech) ([]DefinitionEntry, error)
}

type CardRepo interface {
	Create(context.Context, entity.Card) error
}

type UseCase interface {
	Create(context.Context, string, entity.PartOfSpeech) (entity.Card, error)
}

type DefinitionEntry struct {
	Definition   string
	Provider     string
	PartOfSpeech entity.PartOfSpeech
	Exemples     []string
	Domains      []string
	Registers    []string
}
