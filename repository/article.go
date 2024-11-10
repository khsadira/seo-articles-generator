package repository

import "github.com/qantai/domain"

type Article struct {
	agent domain.Agent
}

func NewArticle(agent domain.Agent) Article {
	return Article{
		agent: agent,
	}
}

func (a Article) GenerateArticle(keyword string) (domain.Article, error) {
	return domain.Article{
		Title:   "title",
		Content: "content" + keyword,
		Status:  "draft",
	}, nil
}
