package entity

import "strings"

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

func ParsePartOfSpeech(pos string) PartOfSpeech {
	switch strings.ToLower(pos) {
	case "noun":
		return Noun
	case "pronoun":
		return Pronoun
	case "verb":
		return Verb
	case "adjective":
		return Adjective
	case "adverb":
		return Adverb
	case "preposition":
		return Preposition
	case "conjunction":
		return Conjunction
	case "interjection":
		return Interjection
	}

	return 0
}
