//nolint:lll
package card

var definitionTemplate = `
----- Vocabulary Card -----
Title: {{ .Title }}
Part Of Speech: {{ printf "%s" .PartOfSpeech }}
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
Part Of Speech: {{ print .PartOfSpeech }}
Definition: {{ .Definition }}
Translations: {{join .Translations ", "}}
---------------------------
`

var definitionPropsTemplate = `
{{- range $index, $element := . -}}
[{{add $index 1}}]{{ if .PartOfSpeech }} {{ .PartOfSpeech }}{{end}}{{ if .Registers }}; {{join .Registers ", "}}{{end}}{{if or .PartOfSpeech .Registers}}{{ if .Domains }}; {{join .Domains ", "}}{{end}} —{{end}} {{ .Definition }}
{{end}}`

var cardListItem = `
{{- printMax .Title 17 }}	{{ printMax (print .PartOfSpeech) 9 }}	{{ printMax .Definition 47 }}	{{ printMax (join .Translations ", ") 20 -}}
`
