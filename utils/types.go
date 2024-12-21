package utils

type Link struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

type Datasource struct {
	Type string `json:"type"`
	UID  string `json:"uid"`
}

type Target struct {
	Expr       string     `json:"expr"`
	Datasource Datasource `json:"datasource"`
}

type RowPanel struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Type        string   `json:"type"`
	Targets     []Target `json:"targets"`
	Panels      []Panel  `json:"panels"`
}

type Panel struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Type        string   `json:"type"`
	Targets     []Target `json:"targets"`
}

type Dashboard struct {
	Title  string     `json:"title"`
	Links  []Link     `json:"links"`
	Panels []RowPanel `json:"panels"`
}

// GetPanels returns all the panels from the dashboard regardless of the nesting
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

// GetPanel returns the panel from the row panel
func (r *RowPanel) GetPanel() Panel {
	var panel Panel
	panel.Title = r.Title
	panel.Description = r.Description
	panel.Type = r.Type
	panel.Targets = r.Targets

	return panel
}
