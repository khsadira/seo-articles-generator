package repository

import "github.com/qantai/domain"

type ArticleRepository struct {
	agent domain.Agent
}

func NewArticle(agent domain.Agent) ArticleRepository {
	return ArticleRepository{
		agent: agent,
	}
}

func (a ArticleRepository) GenerateArticle(keyword, articlePrompt string) (domain.Article, error) {
	return domain.Article{
		Title:   "title_" + keyword,
		Content: "mock_generated: content_" + keyword + " with prompt: " + articlePrompt + "_" + a.agent.ID + "_" + a.agent.APIKey,
		Status:  "draft",
	}, nil
}
