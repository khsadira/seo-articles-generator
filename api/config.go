package api

import "github.com/qantai/domain"

type CMS struct {
	ID     string `json:"ID"`
	URL    string `json:"url"`
	APIKey string `json:"apiKey,omitempty"`
}

type Agent struct {
	ID     string `json:"ID"`
	APIKey string `json:"apiKey"`
}

type ArticlesPublisherConfig struct {
	Keywords      []string `json:"keywords"`
	CMS           []CMS    `json:"cms"`
	PruningAgent  Agent    `json:"pruningAgent"`
	PruningPrompt string   `json:"pruningPrompt"`
	ArticleAgent  Agent    `json:"articleAgent"`
	ArticlePrompt string   `json:"articlePrompt"`
}

type ArticlesPublisherKeywordsConfig struct {
	Keywords      []string `json:"keywords"`
	CMS           []CMS    `json:"cms"`
	ArticleAgent  Agent    `json:"articleAgent"`
	ArticlePrompt string   `json:"articlePrompt"`
}

func toDomainCMS(cms []CMS) []domain.CMS {
	domainCMS := make([]domain.CMS, 0)
	for _, cms := range cms {
		domainCMS = append(domainCMS, domain.CMS{
			ID:     cms.ID,
			URL:    cms.URL,
			APIKey: cms.APIKey,
		})
	}

	return domainCMS
}

func toDomainAgent(agent Agent) domain.Agent {
	return domain.Agent{
		ID:     agent.ID,
		APIKey: agent.APIKey,
	}
}
