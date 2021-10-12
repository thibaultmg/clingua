package language

import (
	"context"
	"fmt"
	"strings"

	"github.com/thibaultmg/clingua/internal/entity"
)

type Dictionnary interface {
	GetDefinition(context.Context, string, entity.PartOfSpeech) ([]DefinitionEntry, error)
}

type Translator interface {
	GetTranslation(context.Context, string, entity.PartOfSpeech) ([]string, error)
}

type LanguageUC interface {
	GetDefinition(context.Context, string, entity.PartOfSpeech) ([]DefinitionEntry, error)
	GetTranslation(context.Context, string, entity.PartOfSpeech) ([]string, error)
}

type DefinitionEntry struct {
	Definition   string
	Provider     string
	PartOfSpeech entity.PartOfSpeech
	Exemples     []string
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
