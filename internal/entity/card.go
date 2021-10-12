package entity

import (
	"golang.org/x/text/language"
)

type Exemple struct {
	Sentence    string
	Translation string
}

type Card struct {
	From         language.Tag // Learner's language
	To           language.Tag // Learned language
	ID           string
	Title        string
	PartOfSpeech PartOfSpeech
	Definition   string
	Exemples     []Exemple
	Translations []string
}

func NewCard() Card {
	return Card{
		From:         language.French,
		To:           language.English,
		ID:           makePseudoUUID(),
		Translations: make([]string, 0, 3),
		Exemples:     make([]Exemple, 0, 2),
	}
}
