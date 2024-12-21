package main

import (
	"log"
	"text/template"
)

func getTemplate() *template.Template {
	mdTemplate := `# {{.Title}}

| Panel Name | Panel Description | Panel Type | Metrics Used |
| ---------- | ----------------- | ---------- | -------- |
{{- range .Panels}}
| {{.Title}} | {{.Description}} | {{.Type}} | {{- range .Metrics}} <pre><code>{{.}}<br></code></pre> {{- end}} |
{{- end}}
`
	tmpl, err := template.New("markdown").Parse(mdTemplate)
	if err != nil {
		log.Fatal(err)
	}

	return tmpl
}
