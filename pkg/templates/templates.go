// Package templates provides markdown template functionality for generating
// documentation from Grafana dashboard data. It contains predefined templates
// and utilities for rendering structured dashboard information into markdown format.
package templates

import (
	"fmt"
	"log/slog"
	"text/template"
)

var (
	// mdTemplate contains the Go template string for generating markdown documentation
	// from Grafana dashboard data. It creates a structured table format with:
	//   - Dashboard title and description as headers
	//   - A table containing panel information with columns for:
	//     * Panel Name
	//     * Panel Description
	//     * Panel Type
	//     * Metrics Used (formatted as inline code blocks)
	//
	// The template uses Go template syntax with range loops to iterate over
	// panels and their associated metrics.
	mdTemplate = `# {{.Title}}
{{.Description}}

| Panel Name | Panel Description | Panel Type | Metrics Used |
| ---------- | ----------------- | ---------- | -------- |
{{- range .Panels}}
| {{.Title}} | {{.Description}} | {{.Type}} | {{- range .Metrics}} ` + "`{{.}}`" + `<br> {{- end}} |
{{- end}}`
)

// GetTemplate creates and returns a parsed Go template for generating markdown
// documentation from dashboard data. The template is based on the predefined
// mdTemplate string and is ready for execution with MarkdownData.
//
// Returns:
//   - *template.Template: A parsed template ready for execution
//   - error: An error if template parsing fails
//
// The returned template expects data conforming to the MarkdownData structure
// from the parser package, containing Title, Description, and Panels fields.
func GetTemplate() (*template.Template, error) {
	tmpl, err := template.New("markdown").Parse(mdTemplate)
	if err != nil {
		slog.Error("error generating a new mardown gotmpl", slog.Any("error", err))
		return nil, fmt.Errorf("error generating a new mardown gotmpl: %w", err)
	}

	return tmpl, nil
}
