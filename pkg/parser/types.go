package parser

// Link represents a dashboard link with its metadata including type, title, and URL.
type Link struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

// Datasource represents a Grafana datasource configuration with type and unique identifier.
type Datasource struct {
	Type string `json:"type"`
	UID  string `json:"uid"`
}

// Target represents a query target containing a PromQL expression and its associated datasource.
type Target struct {
	Expr       string     `json:"expr"`
	Datasource Datasource `json:"datasource"`
}

// RowPanel represents a dashboard row panel that can contain other panels.
// It includes metadata and can hold nested panels within it.
type RowPanel struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Type        string   `json:"type"`
	Targets     []Target `json:"targets"`
	Panels      []Panel  `json:"panels"`
}

// Panel represents a standard dashboard panel with its metadata and query targets.
type Panel struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Type        string   `json:"type"`
	Targets     []Target `json:"targets"`
}

// Dashboard represents a complete Grafana dashboard with its metadata, links, and panels.
type Dashboard struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Links       []Link     `json:"links"`
	Panels      []RowPanel `json:"panels"`
}

// GetPanels returns all panels from the dashboard, flattening the hierarchy
// by extracting panels from row panels and converting row panels themselves
// to regular panels. This provides a unified view of all dashboard content.
func (d *Dashboard) GetPanels() []Panel {
	var panels []Panel
	for _, rowPanel := range d.Panels {
		panels = append(panels, rowPanel.Panels...)
	}
	for _, panel := range d.Panels {
		panels = append(panels, panel.GetPanel())
	}

	return panels
}

// GetPanel converts a RowPanel to a regular Panel by copying its metadata.
// This allows row panels to be treated uniformly with other panel types
// during documentation generation.
func (r *RowPanel) GetPanel() Panel {
	var panel Panel
	panel.Title = r.Title
	panel.Description = r.Description
	panel.Type = r.Type
	panel.Targets = r.Targets

	return panel
}
