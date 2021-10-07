package card

import (
	"context"
	"log"

	"github.com/thibaultmg/clingua/internal/entity"
)

type UseCaseImpl struct {
	cardRepo CardRepo
	defRepo  DefinitionRepo
}

func NewUseCase(cardRepo CardRepo, defRepo DefinitionRepo) UseCase {
	return UseCaseImpl{
		cardRepo: cardRepo,
		defRepo:  defRepo,
	}
}

func (u UseCaseImpl) Create(ctx context.Context, title string, pos entity.PartOfSpeech) (entity.Card, error) {
	ret := entity.NewCard()

	// Get definition
	definition, err := u.defRepo.Get(ctx, title, pos)
	if err != nil {
		log.Panicln(err)
	}

	ret.Definition = definition[0].Definition

	return ret, nil
}
