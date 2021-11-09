package entity

import (
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidPartOfSpeech = errors.New("invalid Part Of Speech")

type PartOfSpeech int

const (
	Any PartOfSpeech = iota
	Noun
	Pronoun
	Verb
	Adjective
	Adverb
	Preposition
	Conjunction
	Interjection
)

func (p PartOfSpeech) IsAny() bool {
	return p == Any
}

func (p PartOfSpeech) String() string {
	switch p {
	case Noun:
		return "noun"
	case Pronoun:
		return "pronoun"
	case Verb:
		return "verb"
	case Adjective:
		return "adjective"
	case Adverb:
		return "adverb"
	case Preposition:
		return "preposition"
	case Conjunction:
		return "conjunction"
	case Interjection:
		return "interjection"
	}

	return ""
}

func ParsePartOfSpeech(pos string) (PartOfSpeech, error) {
	switch strings.ToLower(pos) {
	case "noun":
		return Noun, nil
	case "pronoun":
		return Pronoun, nil
	case "verb":
		return Verb, nil
	case "adjective":
		return Adjective, nil
	case "adverb":
		return Adverb, nil
	case "preposition":
		return Preposition, nil
	case "conjunction":
		return Conjunction, nil
	case "interjection":
		return Interjection, nil
	}

	return 0, fmt.Errorf("error with value %s: %w", pos, ErrInvalidPartOfSpeech)
}
