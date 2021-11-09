package entity

import (
	"golang.org/x/text/language"
)

const maxTranslations = 3

type Card struct {
	From                 language.Tag // Learner's language
	To                   language.Tag // Learned language
	ID                   string
	Title                string
	PartOfSpeech         PartOfSpeech
	Register             string
	Definition           string
	Examples             []string
	ExamplesTranslations []string
	Translations         []string
}

func NewCard() Card {
	return Card{
		From:         language.French,
		To:           language.English,
		ID:           makePseudoUUID(),
		Translations: make([]string, 0, maxTranslations),
		// Examples:     make([]Example, 0, 2),
	}
}
