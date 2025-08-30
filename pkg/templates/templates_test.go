package templates

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTemplate(t *testing.T) {
	tests := []struct {
		name        string
		expectError bool
		description string
		mdTemplate  string
	}{
		{
			name:        "valid template parsing. should return template without error",
			expectError: false,
			description: "Should successfully parse the markdown template with all expected fields",
		},
		{
			name:        "invalid template parsing. should return error",
			expectError: true,
			description: "Should return error when template has invalid syntax",
			mdTemplate: `# {{.Title}}
{{.Description}}

| Panel Name | Panel Description | Panel Type | Metrics Used |
| ---------- | ----------------- | ---------- | -------- |
{{- range .Panels}}
| {{.Title}} | {{.Description}} | {{.Type}} | {{- range .Metrics}} ` + "`{{.}}`" + `<br> {{- end |
{{- end}}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.mdTemplate != "" {
				mdTemplate = tc.mdTemplate
			}
			tmpl, err := GetTemplate()

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, tmpl)
				assert.Contains(t, err.Error(), "error generating a new mardown gotmpl")
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tmpl)

				assert.Equal(t, "markdown", tmpl.Name())
				type Panel struct {
					Title       string
					Description string
					Type        string
					Metrics     []string
				}

				type TemplateData struct {
					Title       string
					Description string
					Panels      []Panel
				}

				testData := TemplateData{
					Title:       "Test",
					Description: "Test Description",
					Panels: []Panel{
						{
							Title:       "Panel1",
							Description: "Desc1",
							Type:        "graph",
							Metrics:     []string{"metric1"},
						},
					},
				}

				var result strings.Builder
				err = tmpl.Execute(&result, testData)
				assert.NoError(t, err, "Template should execute without errors")

				output := result.String()
				assert.Contains(t, output, "# Test")
				assert.Contains(t, output, "Test Description")
				assert.Contains(t, output, "Panel1")
				assert.Contains(t, output, "Desc1")
				assert.Contains(t, output, "graph")
				assert.Contains(t, output, "`metric1`")
			}
		})
	}
}
