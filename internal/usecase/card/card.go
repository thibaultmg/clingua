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

func New(cardRepo CardRepo) CardUCImpl {
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

func (u CardUCImpl) List(ctx context.Context) ([]entity.Card, error) {
	cardsList, err := u.cardRepo.List(ctx)
	if err != nil {
		return []entity.Card{}, fmt.Errorf("unable to list cards from repository: %w", err)
	}

	// TODO: check validity

	return cardsList, nil
}

func (u CardUCImpl) Delete(ctx context.Context, id string) error {
	err := u.cardRepo.Delete(ctx, id)
	if err != nil {
		log.Error().Err(err).Msgf("failed to delete card with ID %s", id)

		return err
	}

	return nil
}
