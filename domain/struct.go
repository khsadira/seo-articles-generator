package domain

const (
	PlaceHolderTypeImage = iota
)

type Agent struct {
	ID          string
	APIKey      string
	Temperature float64
	Model       string
	MaxToken    int
}

type ImageAgent struct {
	ID      string
	APIKey  string
	Size    string
	Quality string
	N       int
}

type PlaceHolder struct {
	ID   string
	Type int
}
