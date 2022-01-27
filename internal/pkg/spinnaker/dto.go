package spinnaker

type DTO struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Text    string `json:"text"`
	Company string `json:"company"`
	PDFURL  string
}

type App struct {
	Name string `json:"name"`
}

type Pipeline struct {
	Name        string            `json:"name"`
	Application string            `json:"application"`
	ID          string            `json:"id"`
	Stages      []Stage           `json:"stages"`
	Variables   map[string]string `json:"variables"`
	Disabled    bool              `json:"disabled"`
	Type        string            `json:"type"`
}

type Stage struct {
	InputArtifacts []InputArtifact `json:"inputArtifacts"`
}

type InputArtifact struct {
	Account  string   `json:"account"`
	Artifact Artifact `json:"artifact"`
}

type Artifact struct {
	ArtfactAccount string `json:"artifactAccount"`
	Id             string `json:"id"`
	Name           string `json:"name"`
	Reference      string `json:"reference"`
	Type           string `json:"type"`
	Version        string `json:"version"`
}
