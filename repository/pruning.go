package repository

import "github.com/seo-articles-generator/domain"

type PruningRepository struct {
	agent domain.Agent
}

func NewPruning(agent domain.Agent) PruningRepository {
	return PruningRepository{
		agent: agent,
	}
}

func (p PruningRepository) GetPrunedKeywords(keyword, pruningPromt string) ([]string, error) {
	return []string{keyword}, nil
}
