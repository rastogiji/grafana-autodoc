package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/rastogiji/autodoc-grafana/utils"
)

type MarkdownData struct {
	Title       string
	Description string
	Panels      []panelData
}

type panelData struct {
	Title       string
	Description string
	Type        string
	Metrics     []string
}

func main() {
	if len(os.Args) != 2 {
		log.Println("A Dashboard file is required as an argument.")
		os.Exit(1)
	}
	dashboard := os.Args[1]
	fmt.Println(dashboard)
	if dashboard == "" || !utils.IsValidFile(dashboard) || !strings.HasSuffix(strings.ToLower(dashboard), ".json") {
		log.Println("Valid dashboard file path is required as an argument.")
		os.Exit(1)
	}

	bs, err := os.ReadFile(dashboard)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.OpenFile(fmt.Sprintf("%s.md", strings.TrimSuffix(filepath.Base(dashboard), ".json")), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var dash utils.Dashboard
	if err := json.Unmarshal(bs, &dash); err != nil {
		log.Fatal(err)
	}

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
				tg := strings.ReplaceAll(target.Expr, "$__range", "1m")
				tg = strings.ReplaceAll(tg, "$__rate_interval", "1m")
				tg = strings.ReplaceAll(tg, "$interval", "1m")
				allMetrics := utils.ExtractMetricFromExpression(tg)
				metrics = append(metrics, allMetrics...)
			}
			newMetrics := utils.GetUniqueElements(metrics)
			pd.Metrics = newMetrics
			data.Panels = append(data.Panels, pd)
		}

	}
	tmpl := getTemplate()
	err = tmpl.Execute(f, data)
	if err != nil {
		log.Fatal(err)
	}

}
