package main

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectError    bool
		errorMsg       string
		expectedInput  string
		expectedOutput string
		expectedLevel  int
		expectedHelp   bool
	}{
		{
			name:           "help flag should return no error and set help to true",
			args:           []string{"program", "--help"},
			expectError:    false,
			expectedHelp:   true,
			expectedInput:  "",
			expectedOutput: ".",
			expectedLevel:  0,
		},
		{
			name:        "missing input flag should return error",
			args:        []string{"program"},
			expectError: true,
			errorMsg:    "input flag is required",
		},
		{
			name:        "invalid log level should return error",
			args:        []string{"program", "--input", "test.json", "--log-level", "99"},
			expectError: true,
			errorMsg:    "invalid log level: 99",
		},
		{
			name:           "valid flags should parse correctly with defaults",
			args:           []string{"program", "--input", "dashboard.json"},
			expectError:    false,
			expectedInput:  "dashboard.json",
			expectedOutput: ".",
			expectedLevel:  0,
			expectedHelp:   false,
		},
		{
			name:           "custom output directory should be parsed correctly",
			args:           []string{"program", "--input", "test.json", "--output", "/custom/output"},
			expectError:    false,
			expectedInput:  "test.json",
			expectedOutput: "/custom/output",
			expectedLevel:  0,
			expectedHelp:   false,
		},
		{
			name:           "custom log level debug should be parsed correctly",
			args:           []string{"program", "--input", "test.json", "--log-level", "-4"},
			expectError:    false,
			expectedInput:  "test.json",
			expectedOutput: ".",
			expectedLevel:  -4,
			expectedHelp:   false,
		},
		{
			name:           "custom log level info should be parsed correctly",
			args:           []string{"program", "--input", "test.json", "--log-level", "0"},
			expectError:    false,
			expectedInput:  "test.json",
			expectedOutput: ".",
			expectedLevel:  0,
			expectedHelp:   false,
		},
		{
			name:           "all flags together should be parsed correctly",
			args:           []string{"program", "--input", "dashboard.json", "--output", "/tmp/docs", "--log-level", "-4"},
			expectError:    false,
			expectedInput:  "dashboard.json",
			expectedOutput: "/tmp/docs",
			expectedLevel:  -4,
			expectedHelp:   false,
		},
		{
			name:           "glob pattern input should be parsed correctly",
			args:           []string{"program", "--input", "*.json", "--output", "./output"},
			expectError:    false,
			expectedInput:  "*.json",
			expectedOutput: "./output",
			expectedLevel:  0,
			expectedHelp:   false,
		},
		{
			name:           "directory input should be parsed correctly",
			args:           []string{"program", "--input", "./dashboards", "--output", "./docs"},
			expectError:    false,
			expectedInput:  "./dashboards",
			expectedOutput: "./docs",
			expectedLevel:  0,
			expectedHelp:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			input = ""
			output = "."
			logLevel = 0
			help = false

			var buf bytes.Buffer
			out = &buf

			os.Args = tc.args

			mockFileProcessor := func() error {
				return nil
			}

			runnerInstance := &runner{
				fileProcessor: mockFileProcessor,
			}

			err := runnerInstance.run()

			if tc.expectError {
				assert.Error(t, err)
				if tc.errorMsg != "" {
					assert.Contains(t, err.Error(), tc.errorMsg)
				}
				return
			}

			assert.NoError(t, err)

			assert.Equal(t, tc.expectedInput, input, "Input flag should be parsed correctly")
			assert.Equal(t, tc.expectedOutput, output, "Output flag should be parsed correctly")
			assert.Equal(t, tc.expectedLevel, logLevel, "Log level flag should be parsed correctly")
			assert.Equal(t, tc.expectedHelp, help, "Help flag should be parsed correctly")

			if !tc.expectedHelp {
				logOutput := buf.String()

				assert.NotEmpty(t, logOutput, "Logger should have been initialized and used")

				assert.Contains(t, logOutput, fmt.Sprintf(`"log-level":%d`, tc.expectedLevel))
				assert.Contains(t, logOutput, fmt.Sprintf(`"input":"%s"`, tc.expectedInput))
				assert.Contains(t, logOutput, fmt.Sprintf(`"output":"%s"`, tc.expectedOutput))
			}
		})
	}
}

func TestProcessFiles(t *testing.T) {
	tests := []struct {
		name         string
		expectError  bool
		input        string
		output       string
		errorMessage string
		setupFiles   func(t *testing.T) string // Returns tmpDir
	}{
		{
			name:        "valid single JSON file should process successfully",
			expectError: false,
			input:       "test.json",
			output:      "output",
			setupFiles: func(t *testing.T) string {
				tmpDir := t.TempDir()
				testFile := filepath.Join(tmpDir, "test.json")
				outputDir := filepath.Join(tmpDir, "output")

				jsonContent := `{
					"dashboard": {
						"title": "Test Dashboard",
						"panels": []
					}
				}`

				err := os.WriteFile(testFile, []byte(jsonContent), 0644)
				assert.NoError(t, err)

				err = os.MkdirAll(outputDir, 0755)
				assert.NoError(t, err)

				return tmpDir
			},
		},
		{
			name:         "non-JSON file should return error",
			expectError:  true,
			input:        "test.txt",
			output:       "output",
			errorMessage: "input file must be a json file",
			setupFiles: func(t *testing.T) string {
				tmpDir := t.TempDir()
				testFile := filepath.Join(tmpDir, "test.txt")
				outputDir := filepath.Join(tmpDir, "output")

				err := os.WriteFile(testFile, []byte("not json"), 0644)
				assert.NoError(t, err)

				err = os.MkdirAll(outputDir, 0755)
				assert.NoError(t, err)

				return tmpDir
			},
		},
		{
			name:        "valid directory with JSON files should process successfully",
			expectError: false,
			input:       "testdir",
			output:      "output",
			setupFiles: func(t *testing.T) string {
				tmpDir := t.TempDir()
				testDir := filepath.Join(tmpDir, "testdir")
				outputDir := filepath.Join(tmpDir, "output")

				err := os.MkdirAll(testDir, 0755)
				assert.NoError(t, err)

				err = os.MkdirAll(outputDir, 0755)
				assert.NoError(t, err)

				// Create multiple JSON files
				for i := 1; i <= 2; i++ {
					jsonContent := `{
						"dashboard": {
							"title": "Test Dashboard ` + string(rune(i+'0')) + `",
							"panels": []
						}
					}`
					testFile := filepath.Join(testDir, "dashboard"+string(rune(i+'0'))+".json")
					err := os.WriteFile(testFile, []byte(jsonContent), 0644)
					assert.NoError(t, err)
				}

				return tmpDir
			},
		},
		{
			name:        "directory with no JSON files should return no error with warning",
			expectError: false,
			input:       "emptydir",
			output:      "output",
			setupFiles: func(t *testing.T) string {
				tmpDir := t.TempDir()
				testDir := filepath.Join(tmpDir, "emptydir")
				outputDir := filepath.Join(tmpDir, "output")

				err := os.MkdirAll(testDir, 0755)
				assert.NoError(t, err)

				err = os.MkdirAll(outputDir, 0755)
				assert.NoError(t, err)

				return tmpDir
			},
		},
		{
			name:        "glob pattern with matching files should process successfully",
			expectError: false,
			input:       "*.json",
			output:      "output",
			setupFiles: func(t *testing.T) string {
				tmpDir := t.TempDir()
				outputDir := filepath.Join(tmpDir, "output")

				err := os.MkdirAll(outputDir, 0755)
				assert.NoError(t, err)

				// Create JSON files in temp directory
				for i := 1; i <= 2; i++ {
					jsonContent := `{
						"dashboard": {
							"title": "Test Dashboard ` + string(rune(i+'0')) + `",
							"panels": []
						}
					}`
					testFile := filepath.Join(tmpDir, "dashboard"+string(rune(i+'0'))+".json")
					err := os.WriteFile(testFile, []byte(jsonContent), 0644)
					assert.NoError(t, err)
				}

				return tmpDir
			},
		},
		{
			name:        "glob pattern with no matching files should return no error with warning",
			expectError: false,
			input:       "nonexistent*.json",
			output:      "output",
			setupFiles: func(t *testing.T) string {
				tmpDir := t.TempDir()
				outputDir := filepath.Join(tmpDir, "output")

				err := os.MkdirAll(outputDir, 0755)
				assert.NoError(t, err)

				return tmpDir
			},
		}, {
			name:         "invalid glob input",
			expectError:  true,
			input:        "[invalid-glob-pattern",
			output:       "output",
			errorMessage: "syntax error in pattern",
			setupFiles: func(t *testing.T) string {
				tmpDir := t.TempDir()
				outputDir := filepath.Join(tmpDir, "output")

				err := os.MkdirAll(outputDir, 0755)
				assert.NoError(t, err)

				return tmpDir
			},
		},
		{
			name:         "invalid input path should return error",
			expectError:  true,
			input:        "nonexistent",
			output:       "output",
			errorMessage: "input path is not a valid file, directory, or glob pattern",
			setupFiles: func(t *testing.T) string {
				tmpDir := t.TempDir()
				outputDir := filepath.Join(tmpDir, "output")

				err := os.MkdirAll(outputDir, 0755)
				assert.NoError(t, err)

				return tmpDir
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup test environment
			tmpDir := tc.setupFiles(t)
			err := os.Chdir(tmpDir)
			assert.NoError(t, err)

			input = tc.input
			output = tc.output

			err = processFiles()

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

func TestValidateFlagValues(t *testing.T) {
	tests := []struct {
		name        string
		logLevel    int
		input       string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid log level debug and valid input. should return no error",
			logLevel:    -4,
			input:       "dashboard.json",
			expectError: false,
		},
		{
			name:        "valid log level info and valid input. should return no error",
			logLevel:    0,
			input:       "dashboard.json",
			expectError: false,
		},
		{
			name:        "valid log level warn and valid input. should return no error",
			logLevel:    4,
			input:       "dashboard.json",
			expectError: false,
		},
		{
			name:        "valid log level error and valid input. should return no error",
			logLevel:    8,
			input:       "dashboard.json",
			expectError: false,
		},
		{
			name:        "invalid log level positive. should return error",
			logLevel:    5,
			input:       "dashboard.json",
			expectError: true,
			errorMsg:    "invalid log level: 5",
		},
		{
			name:        "invalid log level negative. should return error",
			logLevel:    -1,
			input:       "dashboard.json",
			expectError: true,
			errorMsg:    "invalid log level: -1",
		},
		{
			name:        "valid log level but empty input. should return error",
			logLevel:    0,
			input:       "",
			expectError: true,
			errorMsg:    "input flag is required",
		},
		{
			name:        "invalid log level and empty input. should return log level error first",
			logLevel:    99,
			input:       "",
			expectError: true,
			errorMsg:    "invalid log level: 99",
		},
		{
			name:        "valid log level with whitespace input. should return no error",
			logLevel:    0,
			input:       " dashboard.json ",
			expectError: false,
		},
		{
			name:        "valid log level with path input. should return no error",
			logLevel:    -4,
			input:       "/path/to/dashboard.json",
			expectError: false,
		},
		{
			name:        "valid log level with glob pattern. should return no error",
			logLevel:    4,
			input:       "dashboards/*.json",
			expectError: false,
		}, {
			name:        "valid log level with directory. should return no error",
			logLevel:    4,
			input:       "./utils",
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testLogLevel := tc.logLevel
			testInput := tc.input
			logLevel = testLogLevel
			input = testInput

			var buf bytes.Buffer
			setupLog = slog.New(slog.NewJSONHandler(&buf, nil))

			err := validateFlagValues()

			if tc.expectError {
				assert.Error(t, err)
				if tc.errorMsg != "" {
					assert.Contains(t, err.Error(), tc.errorMsg)
				}
				assert.NotEmpty(t, buf.String(), "Expected error to be logged")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
