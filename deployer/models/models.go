package models

type Metadata struct {
	Name string `json:"name"`
}

type Spec struct {
	Image string `json:"image"`
}

// App CRD spec
type App struct {
	ApiVersion string   `json:"apiVersion"`
	Kind       string   `json:"kind"`
	Metadata   Metadata `json:"metadata"`
	Spec       Spec     `json:"spec"`
}

// Patch request config
type Patch struct {
	Spec Spec `json:"spec"`
}

// List of available image tags
type TagsList struct {
	Name string
	Tags []string
}

// Alarm received from Analyzer, providing name of image with anomaly
type Alarm struct {
	Image string `json:"image"`
}
