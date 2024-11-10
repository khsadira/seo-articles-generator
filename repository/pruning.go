package repository

import "github.com/qantai/domain"

type Pruning struct {
	agent domain.Agent
}

func NewPruning(agent domain.Agent) Pruning {
	return Pruning{
		agent: agent,
	}
}

func (p Pruning) GetPrunedKeywords(keyword string) ([]string, error) {
	return []string{keyword}, nil
}
