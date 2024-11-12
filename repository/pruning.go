package repository

import "github.com/qantai/domain"

type PruningRepository struct {
	agent domain.Agent
}

func NewPruning(agent domain.Agent) PruningRepository {
	return PruningRepository{
		agent: agent,
	}
}

func (p PruningRepository) GetPrunedKeywords(keyword, pruningPromt string) ([]string, error) {
	return []string{keyword + "_" + pruningPromt + "_" + p.agent.ID + "_" + p.agent.APIKey}, nil
}
