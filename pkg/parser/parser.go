// Package parser provides functionality for parsing Grafana dashboard JSON files
// and generating structured markdown documentation. It extracts dashboard metadata,
// panel information, and PromQL metrics to create comprehensive documentation.
package parser

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/prometheus/prometheus/promql/parser"
	"github.com/rastogiji/autodoc-grafana/pkg/templates"
	"github.com/rastogiji/autodoc-grafana/pkg/utils"
)

// MarkdownData represents the structured data used for generating markdown documentation
// from a Grafana dashboard. It contains the dashboard's title, description, and
// all panel information formatted for template processing.
type MarkdownData struct {
	// Title is the dashboard title
	Title string
	// Description is the dashboard description
	Description string
	// Panels contains all panels from the dashboard with their metadata and metrics
	Panels []panelData
}

// panelData represents a single dashboard panel with its associated metadata
// and extracted metrics information for documentation purposes.
type panelData struct {
	// Title is the panel title
	Title string
	// Description is the panel description with newlines escaped
	Description string
	// Type indicates the panel type (e.g., "graph", "stat", "table")
	Type string
	// Metrics contains unique metric names extracted from the panel's PromQL queries
	Metrics []string
}

// metricNameVisitor implements the prometheus parser.Visitor interface
// to extract metric names from PromQL expressions through AST traversal.
type metricNameVisitor struct {
	// metricNames stores the collected metric names during AST traversal
	metricNames []string
}

// Visit implements the parser.Visitor interface to traverse PromQL AST nodes
// and extract metric names from VectorSelector nodes. This method is called
// for each node during AST traversal.
func (v *metricNameVisitor) Visit(node parser.Node, path []parser.Node) (parser.Visitor, error) {
	switch n := node.(type) {
	case *parser.VectorSelector:
		v.metricNames = append(v.metricNames, n.Name)
	}
	return v, nil
}

// CreateDocumentationFromFile processes a Grafana dashboard JSON file and generates
// corresponding markdown documentation. It reads the dashboard file, extracts panel
// information and metrics, and writes the formatted documentation to the output directory.
//
// Parameters:
//   - dashboard: path to the Grafana dashboard JSON file
//   - outputDir: directory where the generated markdown file will be saved
//
// Returns an error if file reading, JSON parsing, metric extraction, or template
// execution fails.
func CreateDocumentationFromFile(dashboard string, outputDir string) error {
	logger := slog.With(
		slog.String("processing-file", dashboard),
	)

	logger.Debug("processing file")

	bs, err := os.ReadFile(dashboard)
	if err != nil {
		logger.Error("error reading json file", slog.Any("error", err))
		return fmt.Errorf("error reading dashboard file: %w", err)
	}

	var dash Dashboard
	if err := json.Unmarshal(bs, &dash); err != nil {
		logger.Error("error unmarshalling dashboard json", slog.Any("error", err))
		return fmt.Errorf("error unmarshalling dashboard json: %w", err)
	}

	fileName := filepath.Join(outputDir, strings.TrimSuffix(filepath.Base(dashboard), ".json")+".md")
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		logger.Error("error opening md file", slog.Any("error", err), slog.String("mardown-file", fileName))
		return fmt.Errorf("error opening the corresponding markdown file: %w", err)
	}
	defer f.Close()

	var data MarkdownData

	data.Title = dash.Title
	data.Description = dash.Description

	for _, panel := range dash.GetPanels() {
		var pd panelData
		var metrics []string
		if panel.Type != "row" {
			pd.Title = panel.Title
			pd.Description = strings.ReplaceAll(panel.Description, "\n", "\\n")
			pd.Type = panel.Type
			for _, target := range panel.Targets {
				replacer := strings.NewReplacer("$__range", "1m", "$__rate_interval", "1m", "$interval", "1m")
				tg := replacer.Replace(target.Expr)
				allMetrics, err := extractMetricFromExpression(tg)
				if err != nil {
					return err
				}
				metrics = append(metrics, allMetrics...)
			}
			uniqueMetrics := utils.GetUniqueElements(metrics)
			pd.Metrics = uniqueMetrics
			data.Panels = append(data.Panels, pd)
		}

	}
	tmpl, err := templates.GetTemplate()
	if err != nil {
		return err
	}
	err = tmpl.Execute(f, data)
	if err != nil {
		logger.Error("error executing markdown template", slog.Any("error", err))
		return fmt.Errorf("error executing markdown template: %w", err)
	}
	return nil
}

// extractMetricFromExpression parses a PromQL expression and extracts all metric names
// from it. It handles expression parsing errors and returns the list of unique metrics
// found in the expression.
//
// Parameters:
//   - expr: the PromQL expression string to parse
//
// Returns a slice of metric names and an error if parsing fails.
func extractMetricFromExpression(expr string) ([]string, error) {
	p, err := parser.ParseExpr(expr)
	if err != nil {
		slog.Error("error parsing promql expression", slog.Any("error", err), slog.Any("expr", expr))
		return nil, fmt.Errorf("error parsing promql expression: %w", err)
	}
	return extractMetrics(p), nil
}

// extractMetrics traverses a PromQL AST node and extracts all metric names
// using the metricNameVisitor. This function initiates the AST walk and
// returns the collected metric names.
//
// Parameters:
//   - node: the root PromQL AST node to traverse
//
// Returns a slice of metric names found in the AST.
func extractMetrics(node parser.Node) []string {
	v := &metricNameVisitor{}
	parser.Walk(v, node, nil)
	return v.metricNames
}
