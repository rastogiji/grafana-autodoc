// Package utils provides utility functions for file system operations,
// error handling, and data processing used throughout the grafana-autodoc application.
package utils

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-multierror"
)

// IsValidFile checks if the given file path exists and points to a regular file (not a directory).
// It returns true if the path exists and is a file, false otherwise.
//
// Parameters:
//   - filePath: the path to check
//
// Returns:
//   - true if the path exists and is a regular file
//   - false if the path doesn't exist, is a directory, or any other error occurs
//
// Example:
//
//	if IsValidFile("config.json") {
//	    // File exists and is a regular file
//	}
func IsValidFile(filePath string) bool {
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !fileInfo.IsDir()
}

// IsValidDirectory checks if the given directory path exists.
// It returns true if the path exists (regardless of whether it's a file or directory),
// false if the path doesn't exist.
//
// Note: This function will return true for both files and directories that exist.
// Use IsValidFile if you need to distinguish between files and directories.
//
// Parameters:
//   - dirPath: the directory path to check
//
// Returns:
//   - true if the path exists
//   - false if the path doesn't exist
//
// Example:
//
//	if IsValidDirectory("./dashboards") {
//	    // Path exists
//	}
func IsValidDirectory(dirPath string) bool {
	dirInfo, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return false
	}
	return dirInfo.IsDir()
}

// IsGlobPattern checks if the input string contains glob pattern characters.
// It returns true if the string contains any of the following characters: * ? [ ] { }
// These characters are commonly used in file path globbing patterns.
//
// Parameters:
//   - input: the string to check for glob pattern characters
//
// Returns:
//   - true if the input contains glob pattern characters
//   - false if the input contains no glob pattern characters
//
// Example:
//
//	IsGlobPattern("*.json")        // returns true
//	IsGlobPattern("file.json")     // returns false
//	IsGlobPattern("data/[0-9].txt") // returns true
func IsGlobPattern(input string) bool {
	return strings.ContainsAny(input, "*?[]{}")
}

// SafeMultierrorWait safely waits for a multierror.Group to complete and returns any errors.
// It handles the case where the Group might be nil to prevent panics.
//
// Parameters:
//   - g: pointer to a multierror.Group
//
// Returns:
//   - error: combined errors from the group, or an error if the group is nil
//
// Example:
//
//	var g multierror.Group
//	// ... add goroutines to g ...
//	if err := SafeMultierrorWait(&g); err != nil {
//	    log.Printf("Errors occurred: %v", err)
//	}
func SafeMultierrorWait(g *multierror.Group) error {
	if g == nil {
		return errors.New("multierror.Group is nil")
	}

	err := g.Wait()
	return err.ErrorOrNil()
}

// GetUniqueElements returns a new slice containing only the unique elements from the input slice.
// The order of elements is preserved based on their first occurrence.
// This function works with any comparable type.
//
// Type Parameters:
//   - T: any comparable type (string, int, etc.)
//
// Parameters:
//   - s: slice of comparable elements
//
// Returns:
//   - []T: new slice containing unique elements in order of first occurrence
//
// Example:
//
//	names := []string{"alice", "bob", "alice", "charlie", "bob"}
//	unique := GetUniqueElements(names)
//	// unique = ["alice", "bob", "charlie"]
func GetUniqueElements[T comparable](s []T) []T {
	keys := make(map[T]bool)
	list := []T{}

	for _, entry := range s {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// RetrieveJSONFilesFromDirectory scans a directory and returns a slice of JSON file names.
// It only returns regular files (not directories) that have a .json extension.
// The comparison is case-insensitive, so .JSON, .Json, etc. will also be included.
//
// Parameters:
//   - dirPath: the directory path to scan for JSON files
//
// Returns:
//   - []string: slice of JSON file names (not full paths, just file names)
//   - error: error if the directory cannot be read or doesn't exist
//
// Example:
//
//	files, err := RetrieveJSONFilesFromDirectory("./dashboards")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, file := range files {
//	    fmt.Printf("Found JSON file: %s\n", file)
//	}
func RetrieveJSONFilesFromDirectory(dirPath string) ([]string, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	var jsonFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.ToLower(filepath.Ext(file.Name())) == ".json" {
			jsonFiles = append(jsonFiles, file.Name())
		}
	}
	return jsonFiles, nil
}
