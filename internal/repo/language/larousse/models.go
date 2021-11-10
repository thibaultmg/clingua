package larousse

type responseData struct {
	Data []partOfSpeech
}

type partOfSpeech struct {
	PartOfSpeech string
	Items        []posItem
}

type posItem struct {
	Definition   string
	Translations []string
	Examples     []example
}

type example struct {
	Example     string
	Translation string
}
