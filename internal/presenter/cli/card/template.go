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

var exempleTemplate = `
----- Vocabulary Card -----
Title: {{ .Title }}
{{ if gt (len .Examples) 0 -}}
Examples: {{ index .Examples (getIndex) }}
{{ else -}}
Examples:
{{end}}
---------------------------
`

var exemplesTemplate = `
----- Vocabulary Card -----
Title: {{ .Title }}
Examples: {{ range $index, $element := .Examples}}
	[{{add $index 1}}] — {{ $element }}{{end}}
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
Examples: {{ range $index, $element := .Examples}}
	[{{add $index 1}}] — {{ $element.Example }}; {{ $element.Translation }}{{end}}
---------------------------
`

var definitionPropsTemplate = `
{{- range $index, $element := . -}}
[{{add $index 1}}]{{ if .PartOfSpeech }} {{ .PartOfSpeech }}{{end}}{{ if .Registers }}; {{join .Registers ", "}}{{end}}{{if or .PartOfSpeech .Registers}}{{ if .Domains }}; {{join .Domains ", "}}{{end}} —{{end}} {{ .Definition }}
{{end}}`

var cardListItem = `
{{- printMax .Title 17 }}	{{ printMax (print .PartOfSpeech) 9 }}	{{ printMax .Definition 47 }}	{{ printMax (join .Translations ", ") 20 -}}
`
