package card

var definitionTemplate = `
----- Vocabulary Card -----
Title: {{ .Title }}
Definition: {{ .Definition }}
---------------------------
`

var translationTemplate = `
----- Vocabulary Card -----
Title: {{ .Title }}
Translations: {{join .Translations ", "}}
---------------------------
`

var titleTemplate = `
----- Vocabulary Card -----
Title: {{ .Title }}
---------------------------
`

var cardTemplate = `
----- Vocabulary Card -----
Title: {{ .Title }}
Part Of Speech: {{ .PartOfSpeech }}
Definition: {{ .Definition }}
Translations: {{join .Translations ", "}}
---------------------------
`

var definitionPropsTemplate = `
{{- range $index, $element := . -}}
[{{add $index 1}}]{{ if .PartOfSpeech }} {{ .PartOfSpeech }}{{end}}{{ if .Registers }}; {{join .Registers ", "}}{{end}}{{if or .PartOfSpeech .Registers}}{{ if .Domains }}; {{join .Domains ", "}}{{end}} —{{end}} {{ .Definition }}
{{end}}`
