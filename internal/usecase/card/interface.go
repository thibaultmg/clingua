package card

import (
	"context"

	"github.com/thibaultmg/clingua/internal/entity"
)

type CardRepo interface {
	Get(context.Context, string) (entity.Card, error)
	Create(context.Context, entity.Card) (string, error)
	Delete(context.Context, string) error
	List(context.Context) ([]entity.Card, error)
}

type CardUC interface {
	Create(ctx context.Context, card entity.Card) (string, error)
	List(ctx context.Context) ([]entity.Card, error)
	Delete(context.Context, string) error
}
