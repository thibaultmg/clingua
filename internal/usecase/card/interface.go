package card

import (
	"context"

	"github.com/thibaultmg/clingua/internal/entity"
)

type CardRepo interface {
	Create(context.Context, entity.Card) error
}

type CardUC interface {
	Create(context.Context, string, entity.PartOfSpeech) (entity.Card, error)
}
