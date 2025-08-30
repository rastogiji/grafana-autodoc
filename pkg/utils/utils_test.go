package utils

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"
)

// TestPerson is a comparable struct for testing GetUniqueElements
type TestPerson struct {
	Name string
	Age  int
}

// setupTestEnvironment creates a temporary directory with test files and directories used throughout the tests
func setupTestEnvironment(t *testing.T) (tempDir, testFile, testDir, nestedDir string, cleanup func()) {
	t.Helper()

	tempDir, err := os.MkdirTemp(t.TempDir(), "test_utils")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	testFile = filepath.Join(tempDir, "test.txt")
	file, err := os.Create(testFile)
	if err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	testDir = filepath.Join(tempDir, "testdir")
	err = os.Mkdir(testDir, 0755)
	if err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to create test directory: %v", err)
	}

	nestedDir = filepath.Join(testDir, "nested")
	err = os.Mkdir(nestedDir, 0755)
	if err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to create nested directory: %v", err)
	}

	cleanup = func() {
		os.RemoveAll(tempDir)
	}

	return tempDir, testFile, testDir, nestedDir, cleanup
}

func TestIsValidFile(t *testing.T) {
	tempDir, testFile, testDir, _, cleanup := setupTestEnvironment(t)
	defer cleanup()

	tests := []struct {
		name     string
		filePath string
		expected bool
	}{
		{
			name:     "valid file. should return true.",
			filePath: testFile,
			expected: true,
		}, {
			name:     "glob file type. should return false",
			filePath: filepath.Join(tempDir, "*.json"),
			expected: false,
		},
		{
			name:     "directory type. should resturn false",
			filePath: testDir,
			expected: false,
		},
		{
			name:     "file doesn't exist. should return false",
			filePath: filepath.Join(tempDir, "nonexistent.txt"),
			expected: false,
		},
		{
			name:     "path provided is empty. should return false",
			filePath: "",
			expected: false,
		},
		{
			name:     "invalid path. should return false",
			filePath: "/invalid/path/that/does/not/exist.txt",
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsValidFile(tc.filePath)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestIsValidDirectory(t *testing.T) {
	tempDir, _, testDir, nestedDir, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Test basic directory validation
	t.Run("BasicDirectoryTests", func(t *testing.T) {
		tests := []struct {
			name     string
			dirPath  string
			expected bool
		}{
			{
				name:     "valid directory. should return true",
				dirPath:  testDir,
				expected: true,
			},
			{
				name:     "nested directory. should return true",
				dirPath:  nestedDir,
				expected: true,
			},
			{
				name:     "temp directory itself. should return true",
				dirPath:  tempDir,
				expected: true,
			},
			{
				name:     "non-existent directory. should return false",
				dirPath:  filepath.Join(tempDir, "nonexistent"),
				expected: false,
			},
			{
				name:     "empty path. should return false",
				dirPath:  "",
				expected: false,
			},
			{
				name:     "invalid path. should return false",
				dirPath:  "/invalid/path/that/does/not/exist",
				expected: false,
			},
			{
				name:     "current directory. should return true",
				dirPath:  ".",
				expected: true,
			},
			{
				name:     "parent directory. should return true",
				dirPath:  "..",
				expected: true,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				result := IsValidDirectory(tc.dirPath)
				assert.Equal(t, tc.expected, result)
			})
		}
	})
}

func TestSafeMultierrorWait(t *testing.T) {
	tests := []struct {
		name        string
		setupGroup  func() *multierror.Group
		expectError bool
		errorMsg    string
	}{
		{
			name: "group is nil. should return custom error",
			setupGroup: func() *multierror.Group {
				return nil
			},
			expectError: true,
			errorMsg:    "multierror.Group is nil",
		},
		{
			name: "group is empty i.e no goroutines have been registered",
			setupGroup: func() *multierror.Group {
				var g multierror.Group
				return &g
			},
			expectError: false,
		},
		{
			name: "group has goroutines which all return nil as errors",
			setupGroup: func() *multierror.Group {
				var g multierror.Group
				g.Go(func() error { return nil })
				g.Go(func() error { return nil })
				return &g
			},
			expectError: false,
		},
		{
			name: "all goroutines in the group return error",
			setupGroup: func() *multierror.Group {
				var g multierror.Group
				g.Go(func() error { return errors.New("error 1") })
				g.Go(func() error { return errors.New("error 2") })
				return &g
			},
			expectError: true,
		},
		{
			name: "some goroutines return errors and some don't",
			setupGroup: func() *multierror.Group {
				var g multierror.Group
				g.Go(func() error { return nil })
				g.Go(func() error { return errors.New("only error") })
				g.Go(func() error { return nil })
				return &g
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			group := tc.setupGroup()
			err := SafeMultierrorWait(group)

			if tc.expectError {
				assert.Error(t, err)
				if tc.errorMsg != "" {
					assert.Contains(t, err.Error(), tc.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetUniqueElements(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected any
	}{
		{
			name:     "string slice with duplicates. should return unique elements in order",
			input:    []string{"alice", "bob", "alice", "charlie", "bob", "alice"},
			expected: []string{"alice", "bob", "charlie"},
		},
		{
			name:     "string slice with no duplicates. should return same slice",
			input:    []string{"alice", "bob", "charlie"},
			expected: []string{"alice", "bob", "charlie"},
		},
		{
			name:     "empty string slice. should return empty slice",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "single element string slice. should return same slice",
			input:    []string{"alone"},
			expected: []string{"alone"},
		},
		{
			name:     "int slice with duplicates. should return unique elements in order",
			input:    []int{1, 2, 3, 2, 4, 1, 5, 3},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "int slice with no duplicates. should return same slice",
			input:    []int{1, 2, 3, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "empty int slice. should return empty slice",
			input:    []int{},
			expected: []int{},
		},
		{
			name:     "all same elements. should return single element",
			input:    []string{"same", "same", "same", "same"},
			expected: []string{"same"},
		},
		{
			name: "struct slice with duplicates. should return unique structs in order",
			input: []TestPerson{
				{Name: "Alice", Age: 30},
				{Name: "Bob", Age: 25},
				{Name: "Alice", Age: 30}, // duplicate
				{Name: "Charlie", Age: 35},
				{Name: "Bob", Age: 25}, // duplicate
			},
			expected: []TestPerson{
				{Name: "Alice", Age: 30},
				{Name: "Bob", Age: 25},
				{Name: "Charlie", Age: 35},
			},
		},
		{
			name: "struct slice with no duplicates. should return same slice",
			input: []TestPerson{
				{Name: "Alice", Age: 30},
				{Name: "Bob", Age: 25},
				{Name: "Charlie", Age: 35},
			},
			expected: []TestPerson{
				{Name: "Alice", Age: 30},
				{Name: "Bob", Age: 25},
				{Name: "Charlie", Age: 35},
			},
		},
		{
			name:     "empty struct slice. should return empty slice",
			input:    []TestPerson{},
			expected: []TestPerson{},
		},
		{
			name: "single element struct slice. should return same slice",
			input: []TestPerson{
				{Name: "Alone", Age: 40},
			},
			expected: []TestPerson{
				{Name: "Alone", Age: 40},
			},
		},
		{
			name: "all same struct elements. should return single element",
			input: []TestPerson{
				{Name: "Same", Age: 50},
				{Name: "Same", Age: 50},
				{Name: "Same", Age: 50},
			},
			expected: []TestPerson{
				{Name: "Same", Age: 50},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			switch input := tc.input.(type) {
			case []string:
				result := GetUniqueElements(input)
				expected := tc.expected.([]string)
				assert.Equal(t, expected, result)
			case []int:
				result := GetUniqueElements(input)
				expected := tc.expected.([]int)
				assert.Equal(t, expected, result)
			case []bool:
				result := GetUniqueElements(input)
				expected := tc.expected.([]bool)
				assert.Equal(t, expected, result)
			case []TestPerson:
				result := GetUniqueElements(input)
				expected := tc.expected.([]TestPerson)
				assert.Equal(t, expected, result)
			}
		})
	}
}

func TestIsGlobPattern(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "asterisk wildcard. should return true",
			input:    "*.json",
			expected: true,
		},
		{
			name:     "question mark wildcard. should return true",
			input:    "dashboard?.json",
			expected: true,
		},
		{
			name:     "character class with brackets. should return true",
			input:    "dashboard[0-9].json",
			expected: true,
		},
		{
			name:     "character class with list. should return true",
			input:    "dashboard[abc].json",
			expected: true,
		},
		{
			name:     "brace expansion. should return true",
			input:    "dashboard{1,2,3}.json",
			expected: true,
		},
		{
			name:     "escaped character. should return true",
			input:    "dashboard\\*.json",
			expected: true,
		},
		{
			name:     "multiple wildcards. should return true",
			input:    "**/dashboard*.json",
			expected: true,
		},
		{
			name:     "mixed glob patterns. should return true",
			input:    "test?[0-9]*.json",
			expected: true,
		},
		{
			name:     "regular filename without glob chars. should return false",
			input:    "dashboard.json",
			expected: false,
		},
		{
			name:     "filename with numbers. should return false",
			input:    "dashboard123.json",
			expected: false,
		},
		{
			name:     "filename with spaces. should return false",
			input:    "my dashboard.json",
			expected: false,
		},
		{
			name:     "filename with underscores and dashes. should return false",
			input:    "my_dashboard-v1.json",
			expected: false,
		},
		{
			name:     "path with directory separators. should return false",
			input:    "/path/to/dashboard.json",
			expected: false,
		},
		{
			name:     "relative path. should return false",
			input:    "./dashboard.json",
			expected: false,
		},
		{
			name:     "empty string. should return false",
			input:    "",
			expected: false,
		},
		{
			name:     "just a period. should return false",
			input:    ".",
			expected: false,
		},
		{
			name:     "double period. should return false",
			input:    "..",
			expected: false,
		},
		{
			name:     "filename with extension only. should return false",
			input:    ".json",
			expected: false,
		},
		{
			name:     "complex path with glob. should return true",
			input:    "/path/to/dashboards/*.json",
			expected: true,
		},
		{
			name:     "negation pattern. should return true",
			input:    "dashboard[!0-9].json",
			expected: true,
		},
		{
			name:     "range with uppercase. should return true",
			input:    "dashboard[A-Z].json",
			expected: true,
		},
		{
			name:     "single character in brackets. should return true",
			input:    "dashboard[a].json",
			expected: true,
		},
		{
			name:     "asterisk at beginning. should return true",
			input:    "*dashboard.json",
			expected: true,
		},
		{
			name:     "asterisk in middle. should return true",
			input:    "dash*board.json",
			expected: true,
		},
		{
			name:     "question mark at end. should return true",
			input:    "dashboard?",
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsGlobPattern(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestRetrieveJSONFilesFromDirectory(t *testing.T) {
	tests := []struct {
		name          string
		setupDir      func(t *testing.T) (string, func())
		expectedFiles []string
		expectError   bool
		errorMsg      string
	}{
		{
			name: "directory with json files. should return all json files case insensitive",
			setupDir: func(t *testing.T) (string, func()) {
				tempDir, err := os.MkdirTemp(t.TempDir(), "json_test")
				assert.NoError(t, err)

				files := []string{"dashboard1.json", "config.JSON", "data.Json", "settings.json"}
				for _, file := range files {
					f, err := os.Create(filepath.Join(tempDir, file))
					assert.NoError(t, err)
					f.Close()
				}

				return tempDir, func() { os.RemoveAll(tempDir) }
			},
			expectedFiles: []string{"config.JSON", "dashboard1.json", "data.Json", "settings.json"},
			expectError:   false,
		},
		{
			name: "directory with mixed file types. should return only json files",
			setupDir: func(t *testing.T) (string, func()) {
				tempDir, err := os.MkdirTemp(t.TempDir(), "mixed_test")
				assert.NoError(t, err)

				files := []string{"valid.json", "another.JSON", "readme.txt", "config.yaml", "script.sh"}

				for _, file := range files {
					f, err := os.Create(filepath.Join(tempDir, file))
					assert.NoError(t, err)
					f.Close()
				}

				return tempDir, func() { os.RemoveAll(tempDir) }
			},
			expectedFiles: []string{"another.JSON", "valid.json"},
			expectError:   false,
		},
		{
			name: "directory with subdirectories. should ignore subdirectories and nested files",
			setupDir: func(t *testing.T) (string, func()) {
				tempDir, err := os.MkdirTemp(t.TempDir(), "subdir_test")
				assert.NoError(t, err)

				f, err := os.Create(filepath.Join(tempDir, "root.json"))
				assert.NoError(t, err)
				f.Close()

				subdir := filepath.Join(tempDir, "subdir")
				err = os.Mkdir(subdir, 0755)
				assert.NoError(t, err)

				f, err = os.Create(filepath.Join(subdir, "nested.json"))
				assert.NoError(t, err)
				f.Close()

				return tempDir, func() { os.RemoveAll(tempDir) }
			},
			expectedFiles: []string{"root.json"},
			expectError:   false,
		},
		{
			name: "empty directory. should return empty slice",
			setupDir: func(t *testing.T) (string, func()) {
				tempDir, err := os.MkdirTemp(t.TempDir(), "empty_test")
				assert.NoError(t, err)
				return tempDir, func() { os.RemoveAll(tempDir) }
			},
			expectedFiles: nil,
			expectError:   false,
		},
		{
			name: "non existent directory. should return error",
			setupDir: func(t *testing.T) (string, func()) {
				return "/path/that/does/not/exist", func() {}
			},
			expectedFiles: nil,
			expectError:   true,
			errorMsg:      "no such file or directory",
		},
		{
			name: "path is file not directory. should return error",
			setupDir: func(t *testing.T) (string, func()) {
				tempDir, err := os.MkdirTemp(t.TempDir(), "file_test")
				assert.NoError(t, err)

				filePath := filepath.Join(tempDir, "notadir.txt")
				f, err := os.Create(filePath)
				assert.NoError(t, err)
				f.Close()

				return filePath, func() { os.RemoveAll(tempDir) }
			},
			expectedFiles: nil,
			expectError:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dirPath, cleanup := tc.setupDir(t)
			defer cleanup()

			result, err := RetrieveJSONFilesFromDirectory(dirPath)

			if tc.expectError {
				assert.Error(t, err)
				if tc.errorMsg != "" {
					assert.Contains(t, err.Error(), tc.errorMsg)
				}
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tc.expectedFiles, result)
			}
		})
	}
}
