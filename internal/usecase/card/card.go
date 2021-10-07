package card

import (
	"context"

	"github.com/thibaultmg/clingua/internal/entity"
)

type CardUCImpl struct {
	cardRepo CardRepo
}

func NewUseCase(cardRepo CardRepo) CardUC {
	return CardUCImpl{
		cardRepo: cardRepo,
	}
}

func (u CardUCImpl) Create(ctx context.Context, title string, pos entity.PartOfSpeech) (entity.Card, error) {
	ret := entity.NewCard()

	// // Get definition
	// definition, err := u.defRepo.GetDefinition(ctx, title, pos)
	// if err != nil {
	// 	log.Panicln(err)
	// }

	// ret.Definition = definition[0].Definition

	return ret, nil
}
