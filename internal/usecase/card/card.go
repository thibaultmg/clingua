package card

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/thibaultmg/clingua/internal/entity"
)

type CardUCImpl struct {
	cardRepo CardRepo
}

func New(cardRepo CardRepo) CardUC {
	return CardUCImpl{
		cardRepo: cardRepo,
	}
}

func (u CardUCImpl) Create(ctx context.Context, card entity.Card) (string, error) {
	log.Debug().Msgf("creating card with title %s", card.Title)

	// TODO: check card validity

	cardID, err := u.cardRepo.Create(ctx, card)
	if err != nil {
		return "", fmt.Errorf("unable to create card in repository: %w", err)
	}

	return cardID, nil
}
