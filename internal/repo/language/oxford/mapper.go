package oxford

import (
	"log"
	"strings"

	"github.com/thibaultmg/clingua/internal/entity"
	"github.com/thibaultmg/clingua/internal/usecase/language"
)

func response2Internal(response EntriesResponse) []language.DefinitionEntry {
	var ret []language.DefinitionEntry

	for _, r := range response.Results {
		for _, le := range r.LexicalEntries {
			for _, e := range le.Entries {
				for _, s := range e.Senses {
					pos, err := entity.ParsePartOfSpeech(le.LexicalCategory.Text)
					if err != nil {
						log.Printf("Error with sense %#v: ", s)
					}

					newEntry := language.DefinitionEntry{
						Definition:   s.Definitions[0],
						PartOfSpeech: pos,
					}

					for _, ex := range s.Examples {
						newEntry.Examples = append(newEntry.Examples, ex.Text)
					}

					for _, reg := range s.Registers {
						newEntry.Registers = append(newEntry.Registers, strings.ToLower(reg.Text))
					}

					newEntry.Provider = response.Metadata.Provider

					// Skipping sub senses for now...
					// for _, sub := range s.Subsenses {
					// 	newEntry.Domains = append(newEntry.Domains, sub.Domains)
					// }

					ret = append(ret, newEntry)
				}
			}
		}
	}

	return ret
}
