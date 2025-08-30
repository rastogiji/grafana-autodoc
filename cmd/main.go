// Package main provides a command-line tool for automatically generating
// documentation from Grafana dashboard JSON files. It supports processing
// single files, directories, or glob patterns and outputs structured
// markdown documentation.
package main

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/rastogiji/autodoc-grafana/pkg/parser"
	"github.com/rastogiji/autodoc-grafana/pkg/utils"
	flag "github.com/spf13/pflag"
)

var (
	// input specifies the path to dashboard file, directory, or glob pattern
	input string
	// output specifies the path to output directory where markdown files will be generated
	output string
	// logLevel sets the logging level (Debug: -4, Info: 0, Warn: 4, Error: 8)
	logLevel int
	// help indicates whether to show the help message
	help bool
	// setupLog is the initial logger used for setup and validation phases
	setupLog = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	// out is the output writer, configurable for testing
	out io.Writer = os.Stdout
)

// runner encapsulates the application's main execution logic with
// dependency injection for testability.
type runner struct {
	// fileProcessor is the function that handles the actual file processing logic.
	// It can be injected for testing purposes.
	fileProcessor func() error
}

// main is the entry point of the application. It initializes the runner
// with the default file processor and executes the main application logic.
func main() {
	autodocRunner := runner{
		fileProcessor: processFiles,
	}

	if err := autodocRunner.run(); err != nil {
		os.Exit(1)
	}
}

// run executes the main application logic including command-line flag parsing,
// validation, logger configuration, and file processing. It returns an error
// if any step fails.
func (r *runner) run() error {
	cli := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	cli.StringVar(&input, "input", "", "Path to dashboard file, directory, or glob pattern (e.g., dashboard.json, ./dashboards, files/*.json)")
	cli.StringVar(&output, "output", ".", "Path to output directory where markdown files will be generated (default: current directory)")
	cli.IntVar(&logLevel, "log-level", 0, "Debug: -4, Info: 0, Warn: 4, Error: 8 (default: Info)")
	cli.BoolVar(&help, "help", false, "Show help message")

	cli.Parse(os.Args[1:])

	if help {
		flag.Usage()
		return nil
	}

	if err := validateFlagValues(); err != nil {
		return err
	}

	logger := slog.New(slog.NewJSONHandler(out, &slog.HandlerOptions{
		Level: slog.Level(logLevel),
	})).With(
		slog.Int("log-level", logLevel),
		slog.String("input", input),
		slog.String("output", output),
	)

	slog.SetDefault(logger)

	slog.Info("beginning processing files")
	return r.fileProcessor()
}

// processFiles handles the actual file processing logic based on the input type.
// It supports three input modes:
//   - Glob patterns: processes all matching files
//   - Single files: processes a single JSON file
//   - Directories: processes all JSON files in the directory
//
// Returns an error if processing fails for any file.
func processFiles() error {
	switch {
	case utils.IsGlobPattern(input):
		matches, err := filepath.Glob(input)
		if err != nil {
			slog.Error("Error processing glob pattern", slog.Any("error", err))
			return err
		}
		if len(matches) == 0 {
			slog.Warn("No files found matching pattern")
			return nil
		}
		slog.Info("Found files matching pattern", slog.Int("file-count", len(matches)))
		var g multierror.Group
		for _, match := range matches {
			g.Go(func() error {
				if strings.ToLower(filepath.Ext(match)) == ".json" {
					if err := parser.CreateDocumentationFromFile(match, output); err != nil {
						return err
					}
				} else {
					slog.Debug("Skipping non-JSON file", slog.String("file", match))
				}
				return nil
			})
		}

		if err := utils.SafeMultierrorWait(&g); err != nil {
			slog.Error("error processing files", slog.Any("error", err))
			return err
		}
		slog.Info("Processed JSON files from glob pattern")
		return nil
	case utils.IsValidFile(input):
		if strings.ToLower(filepath.Ext(input)) != ".json" {
			slog.Error("Input file must be a JSON file")
			return errors.New("input file must be a json file")
		}
		slog.Info("Processing single file")
		if err := parser.CreateDocumentationFromFile(input, output); err != nil {
			return err
		}
		return nil
	case utils.IsValidDirectory(input):
		files, err := utils.RetrieveJSONFilesFromDirectory(input)
		if err != nil {
			slog.Error("Error retrieving files from directory", slog.Any("error", err))
			return err
		}

		if len(files) == 0 {
			slog.Warn("No JSON files found in directory")
			return nil
		}

		slog.Info("Found JSON files in directory", slog.Int("count", len(files)))
		var g multierror.Group
		for _, file := range files {
			g.Go(func() error {
				dashboard := filepath.Join(input, file)
				if err := parser.CreateDocumentationFromFile(dashboard, output); err != nil {
					return err
				}
				return nil
			})
		}

		if err := utils.SafeMultierrorWait(&g); err != nil {
			slog.Error("error processing files", slog.Any("error", err))
			return err
		}
		slog.Info("Processed all files in the directory")
		return nil
	default:
		slog.Error("Input path is not a valid file, directory, or glob pattern")
		return errors.New("input path is not a valid file, directory, or glob pattern")
	}
}

// validateFlagValues validates the command-line flag values to ensure they
// meet the application's requirements. It checks that:
//   - logLevel is one of the valid values: -4 (Debug), 0 (Info), 4 (Warn), 8 (Error)
//   - input flag is provided and not empty
//
// Returns an error if validation fails.
func validateFlagValues() error {
	validLogLevels := []int{-4, 0, 4, 8}
	isValid := false
	for _, valid := range validLogLevels {
		if logLevel == valid {
			isValid = true
			break
		}
	}
	if !isValid {
		setupLog.Error("Invalid log level", slog.Int("log-level", logLevel), slog.String("valid_values", "Debug(-4), Info(0), Warn(4), Error(8)"))
		return fmt.Errorf("invalid log level: %d", logLevel)
	}

	if input == "" {
		setupLog.Error("input flag is required")
		return errors.New("input flag is required")
	}
	return nil
}
