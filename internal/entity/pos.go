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
	PhrasalVerb
	Adjective
	Adverb
	Preposition
	Conjunction
	Interjection
	Idiom
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
	case PhrasalVerb:
		return "phrasal verb"
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
	case Idiom:
		return "idiom"
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
	case "phrasal verb":
		return PhrasalVerb, nil
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
	case "idiom":
		return Idiom, nil
	}

	return 0, fmt.Errorf("error with value %s: %w", pos, ErrInvalidPartOfSpeech)
}
