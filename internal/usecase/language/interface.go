package language

import (
	"context"
	"fmt"
	"strings"

	"github.com/thibaultmg/clingua/internal/entity"
)

type Dictionary interface {
	Define(context.Context, string, entity.PartOfSpeech) ([]DefinitionEntry, error)
}

type Translator interface {
	Translate(context.Context, string) ([]string, error)
}

type WordTranslator interface {
	TranslateWord(context.Context, string, entity.PartOfSpeech) ([]WordTranslationEntry, error)
}

type Sentences interface {
	Sentences(context.Context, string) ([]string, error)
}

type LanguageUC interface {
	Dictionary
	Translator
	WordTranslator
	Sentences
}

type WordTranslationEntry struct {
	PartOfSpeech entity.PartOfSpeech
	Translation  string
	Meaning      string
}

type DefinitionEntry struct {
	Definition   string
	Provider     string
	PartOfSpeech entity.PartOfSpeech
	Examples     []string
	Domains      []string
	Registers    []string
}

func (d DefinitionEntry) String() string {
	var ret strings.Builder

	if d.PartOfSpeech != entity.Any {
		// regs := strings.Join(d.Registers, ", ")
		fmt.Fprintf(&ret, "%s ", d.PartOfSpeech.String())
	}

	if len(d.Registers) > 0 {
		regs := strings.Join(d.Registers, ", ")
		fmt.Fprintf(&ret, "(%s) ", regs)
	}

	if len(d.Domains) > 0 {
		regs := strings.Join(d.Domains, ", ")
		fmt.Fprintf(&ret, "(%s) ", regs)
	}

	fmt.Fprintf(&ret, "-- %s", d.Definition)
	// ret.WriteString(d.Definition)

	if len(d.Provider) > 0 {
		fmt.Fprintf(&ret, " -- %s", d.Provider)
	}

	return ret.String()
}
