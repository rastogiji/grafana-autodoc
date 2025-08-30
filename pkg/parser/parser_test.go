package parser

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDocumentationFromFile(t *testing.T) {
	tempOutputDir, err := os.MkdirTemp(t.TempDir(), "parser_test")
	if err != nil {
		t.Fatalf("error creating temp directory: %v", err)
	}

	tests := []struct {
		name         string
		filename     string
		expectError  bool
		outputDir    string
		errorMessage string
	}{
		{
			name:        "valid json file should create the correct markdown file. no errors",
			filename:    "testdata/valid_dashboard.json",
			expectError: false,
		}, {
			name:         "empty file. should return error",
			filename:     "testdata/empty_dashboard.json",
			expectError:  true,
			errorMessage: "error unmarshalling dashboard json: unexpected end of JSON input",
		}, {
			name:         "bad query in panel should return error",
			filename:     "testdata/bad_query.json",
			expectError:  true,
			errorMessage: "error parsing promql expression",
		}, {
			name:         "bad json schema should return error",
			filename:     "testdata/bad_schema.json",
			expectError:  true,
			errorMessage: "error unmarshalling dashboard json",
		}, {
			name:         "json file doesnt exist should return error",
			filename:     "/file/that/does/not/exist",
			expectError:  true,
			errorMessage: "error reading dashboard file",
		}, {
			name:         "error opening the mardown file for docs",
			filename:     "testdata/valid_dashboard.json",
			outputDir:    "/invalid/output/directory/which/does/not/exist",
			expectError:  true,
			errorMessage: "error opening the corresponding markdown file",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var outputDir string
			if tc.outputDir == "" {
				outputDir = tempOutputDir
			} else {
				outputDir = tc.outputDir
			}
			err = CreateDocumentationFromFile(tc.filename, outputDir)
			if tc.expectError {
				assert.Error(t, err)
				if tc.errorMessage != "" {
					assert.Contains(t, err.Error(), tc.errorMessage)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
