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
}

func NewCard() Card {
	return Card{
		From: language.French,
		To:   language.English,
		ID:   makePseudoUUID(),
	}
}
