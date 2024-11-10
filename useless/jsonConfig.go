package useless

type CMS struct {
	ID       string `json:"ID"`
	Auth     string `json:"auth"`
	APIKey   string `json:"apiKey,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type Agent struct {
	ID     string `json:"ID"`
	APIKey string `json:"apiKey"`
}

type Config struct {
	PrunedKeywords []string `json:"prunedKeywords"`
	Keywords       []string `json:"keywords"`
	CMS            []CMS    `json:"cms"`
	PruningAgent   Agent    `json:"pruningAgent"`
	ArticleAgent   Agent    `json:"articleAgent"`
}
