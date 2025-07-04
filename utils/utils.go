package utils

import (
	"fmt"
	"github.com/prometheus/prometheus/promql/parser"
	"os"
	"strings"
)

func IsValidFile(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func IsValidDirectory(dirPath string) bool {
	_, err := os.Stat(dirPath)
	return !os.IsNotExist(err)
}

func ExtractMetricFromExpression(expr string) []string {
	p, err := parser.ParseExpr(expr)
	if err != nil {
		fmt.Printf("Error:%v\nExpression:%s\n", err, expr)
	}
	return extractMetrics(p)
}

type metricNameVisitor struct {
	metricNames []string
}

func (v *metricNameVisitor) Visit(node parser.Node, path []parser.Node) (parser.Visitor, error) {
	switch n := node.(type) {
	case *parser.VectorSelector:
		v.metricNames = append(v.metricNames, n.Name)
	}
	return v, nil
}

func extractMetrics(node parser.Node) []string {
	v := &metricNameVisitor{}
	parser.Walk(v, node, nil)
	return v.metricNames
}

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

func RetrieveFilesFromDirectory(dirPath string) ([]string, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	var jsonFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".json") {
			jsonFiles = append(jsonFiles, file.Name())
		}
	}
	return jsonFiles, nil
}
