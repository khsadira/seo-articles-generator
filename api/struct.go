package api

import "github.com/qantai/domain"

type CMS struct {
	ID     string `json:"ID"`
	URL    string `json:"url"`
	User   string `json:"user"`
	APIKey string `json:"apiKey"`
}

type Agent struct {
	ID          string  `json:"ID"`
	APIKey      string  `json:"apiKey"`
	Model       string  `json:"model"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"maxTokens"`
}

type ImageAgent struct {
	ID      string `json:"ID"`
	APIKey  string `json:"apiKey"`
	Size    string `json:"size"`
	Quality string `json:"quality"`
	N       int    `json:"n"`
}

type ArticlesPublisherConfig struct {
	Keywords      []string `json:"keywords"`
	CMS           []CMS    `json:"cms"`
	PruningAgent  Agent    `json:"pruningAgent"`
	PruningPrompt string   `json:"pruningPrompt"`
	ArticleAgent  Agent    `json:"articleAgent"`
	ArticlePrompt string   `json:"articlePrompt"`
}

type ArticlesPublisherPrunedKeywordsConfig struct {
	Keywords      []string   `json:"keywords"`
	CMS           []CMS      `json:"cms"`
	ArticleAgent  Agent      `json:"articleAgent"`
	ArticlePrompt string     `json:"articlePrompt"`
	ImageAgent    ImageAgent `json:"imageAgent"`
	ImagePrompt   string     `json:"imagePrompt"`
}

func toDomainCMS(cms []CMS) []domain.CMS {
	domainCMS := make([]domain.CMS, 0)
	for _, cms := range cms {
		domainCMS = append(domainCMS, domain.CMS{
			ID:     cms.ID,
			URL:    cms.URL,
			User:   cms.User,
			APIKey: cms.APIKey,
		})
	}

	return domainCMS
}

func toDomainAgent(agent Agent) domain.Agent {
	return domain.Agent{
		ID:          agent.ID,
		APIKey:      agent.APIKey,
		Model:       agent.Model,
		Temperature: agent.Temperature,
		MaxToken:    agent.MaxTokens,
	}
}

func toDomainImageAgent(agent ImageAgent) domain.ImageAgent {
	return domain.ImageAgent{
		ID:      agent.ID,
		APIKey:  agent.APIKey,
		Size:    agent.Size,
		Quality: agent.Quality,
		N:       agent.N,
	}
}
