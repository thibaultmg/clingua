package card

var definitionTemplate = `
----- Vocabulary Card -----
Title: {{ .Title }}
Definition: {{ .Definition }}
---------------------------
`

var titleTemplate = `
----- Vocabulary Card -----
Title: {{ .Title }}
Definition: {{ .Definition }}
---------------------------
`

var definitionPropsTemplate = `
{{- range $index, $element := . -}}
[{{add $index 1}}]{{ if .PartOfSpeech }} {{ .PartOfSpeech }}{{end}}{{ if .Registers }}; {{join .Registers ", "}}{{end}}{{if or .PartOfSpeech .Registers}}{{ if .Domains }}; {{join .Domains ", "}}{{end}} —{{end}} {{ .Definition }}
{{end}}`
