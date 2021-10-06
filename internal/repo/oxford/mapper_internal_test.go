package oxford

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thibaultmg/clingua/internal/entity"
)

var jsonResponse = `{"id":"ace","metadata":{"operation":"retrieve","provider":"Oxford University Press","schema":"RetrieveEntry"},"results":[{"id":"ace","language":"en-gb","lexicalEntries":[{"entries":[{"homographNumber":"102","pronunciations":[{"audioFile":"https://audio.oxforddictionaries.com/en/mp3/ace_1_gb_1_abbr.mp3","dialects":["British English"],"phoneticNotation":"IPA","phoneticSpelling":"eɪs"}],"senses":[{"definitions":["(in tennis and similar games) serve an ace against (an opponent)"],"examples":[{"text":"he can ace opponents with serves of no more than 62 mph"}],"id":"m_en_gbus0005680.020","registers":[{"id":"informal","text":"Informal"}],"subsenses":[{"definitions":["score an ace on (a hole) or with (a shot)"],"domains":[{"id":"golf","text":"Golf"}],"examples":[{"text":"there was a prize for the first player to ace the hole"}],"id":"m_en_gbus0005680.026"}]},{"definitions":["achieve high marks in (a test or exam)"],"examples":[{"text":"I aced my grammar test"}],"id":"m_en_gbus0005680.028","registers":[{"id":"informal","text":"Informal"}],"subsenses":[{"definitions":["outdo someone in a competitive situation"],"examples":[{"text":"the magazine won an award, acing out its rivals"}],"id":"m_en_gbus0005680.029"}]}]}],"language":"en-gb","lexicalCategory":{"id":"verb","text":"Verb"},"text":"ace"}],"type":"headword","word":"ace"}],"word":"ace"}`

func TestOxford_Mapper_Nominal(t *testing.T) {
	assert := assert.New(t)

	var resp EntriesResponse
	err := json.Unmarshal([]byte(jsonResponse), &resp)
	assert.Nil(err)
	internal := response2Internal(resp)
	assert.Len(internal, 2)
	// Check first item
	assert.Equal("(in tennis and similar games) serve an ace against (an opponent)", internal[0].Definition)
	assert.Equal(internal[0].PartOfSpeech, entity.Verb)
	assert.Len(internal[0].Exemples, 1)
	assert.Equal("he can ace opponents with serves of no more than 62 mph", internal[0].Exemples[0])
	assert.Equal("informal", internal[0].Registers[0])

	// Check second item
	assert.Equal("achieve high marks in (a test or exam)", internal[1].Definition)
	assert.Equal(internal[1].PartOfSpeech, entity.Verb)
	assert.Len(internal[1].Exemples, 1)
	assert.Equal("I aced my grammar test", internal[1].Exemples[0])
	assert.Equal("informal", internal[1].Registers[0])
}
