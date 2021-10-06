package oxford

type EntriesResponse struct {
	ID       string    `json:"id"`
	Metadata Metadata  `json:"metadata"`
	Results  []Results `json:"results"`
	Word     string    `json:"word"`
}

type Metadata struct {
	Operation string `json:"operation"`
	Provider  string `json:"provider"`
	Schema    string `json:"schema"`
}
type Pronunciations struct {
	AudioFile        string   `json:"audioFile"`
	Dialects         []string `json:"dialects"`
	PhoneticNotation string   `json:"phoneticNotation"`
	PhoneticSpelling string   `json:"phoneticSpelling"`
}
type Examples struct {
	Text string `json:"text"`
}
type Registers struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}
type Domains struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}
type Subsenses struct {
	Definitions []string   `json:"definitions"`
	Domains     []Domains  `json:"domains"`
	Examples    []Examples `json:"examples"`
	ID          string     `json:"id"`
}
type Senses struct {
	Definitions []string    `json:"definitions"`
	Examples    []Examples  `json:"examples"`
	ID          string      `json:"id"`
	Registers   []Registers `json:"registers"`
	Subsenses   []Subsenses `json:"subsenses"`
}
type Entries struct {
	HomographNumber string           `json:"homographNumber"`
	Pronunciations  []Pronunciations `json:"pronunciations"`
	Senses          []Senses         `json:"senses"`
}
type LexicalCategory struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}
type LexicalEntries struct {
	Entries         []Entries       `json:"entries"`
	Language        string          `json:"language"`
	LexicalCategory LexicalCategory `json:"lexicalCategory"`
	Text            string          `json:"text"`
}
type Results struct {
	ID             string           `json:"id"`
	Language       string           `json:"language"`
	LexicalEntries []LexicalEntries `json:"lexicalEntries"`
	Type           string           `json:"type"`
	Word           string           `json:"word"`
}
