package oxford

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thibaultmg/clingua/internal/entity"
)

//go:embed testdata/response.json
var jsonResponse []byte

func TestOxford_Mapper_Nominal(t *testing.T) {
	assert := assert.New(t)
	t.Parallel()

	var resp EntriesResponse
	err := json.Unmarshal(jsonResponse, &resp)
	assert.Nil(err)

	internal := response2Internal(resp)
	assert.Len(internal, 2)
	// Check first item
	assert.Equal("(in tennis and similar games) serve an ace against (an opponent)", internal[0].Definition)
	assert.Equal(internal[0].PartOfSpeech, entity.Verb)
	assert.Len(internal[0].Examples, 1)
	assert.Equal("he can ace opponents with serves of no more than 62 mph", internal[0].Examples[0])
	assert.Equal("informal", internal[0].Registers[0])
	assert.Equal("Oxford University Press", internal[0].Provider)

	// Check second item
	assert.Equal("achieve high marks in (a test or exam)", internal[1].Definition)
	assert.Equal(internal[1].PartOfSpeech, entity.Verb)
	assert.Len(internal[1].Examples, 1)
	assert.Equal("I aced my grammar test", internal[1].Examples[0])
	assert.Equal("informal", internal[1].Registers[0])
}
