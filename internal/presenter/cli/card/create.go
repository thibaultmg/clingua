package card

import (
	"context"
	"fmt"
	"time"

	"github.com/manifoldco/promptui"

	"github.com/thibaultmg/clingua/internal/entity"
	"github.com/thibaultmg/clingua/internal/usecase/language"
)

type CardPresenter struct {
	languageUC language.LanguageUC
	definition language.DefinitionEntry
	card       entity.Card
}

func New(luc language.LanguageUC) *CardPresenter {
	return &CardPresenter{
		languageUC: luc,
	}
}

func (c *CardPresenter) CreateCard(word string, pos entity.PartOfSpeech) error {
	c.card.Title = word

	// Get definition.
	def, err := c.selectDefinition(word, pos)
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return fmt.Errorf("failed to get definition: %v", err)
	}

	c.card.Definition = def.Definition
	c.card.PartOfSpeech = def.PartOfSpeech

	// cache definition
	c.definition = def

	// Get translation

	// Get exemples

	// Get Word (and exemple?) speech

	return nil
}

func (c CardPresenter) selectDefinition(word string, pos entity.PartOfSpeech) (language.DefinitionEntry, error) {
	var ret language.DefinitionEntry

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	defs, err := c.languageUC.GetDefinition(ctx, word, pos)
	if err != nil {
		return ret, fmt.Errorf("failed to get definitions: %v", err)
	}

	items := make([]string, 0, len(defs))
	for i, v := range defs {
		def := fmt.Sprintf("[%d] %v", i+1, v)
		items = append(items, def)
	}

	prompt := promptui.Select{
		Label: "Select the definition",
		Items: items,
	}

	resultIdx, _, err := prompt.Run()
	if err != nil {
		return ret, fmt.Errorf("failed to get selected definition from prompt: %v", err)
	}

	ret = defs[resultIdx]

	return ret, nil

}
