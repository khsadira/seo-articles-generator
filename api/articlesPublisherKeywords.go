package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/seo-articles-generator/domain"
	"github.com/seo-articles-generator/repository"
)

func HandlerPublishArticles(w http.ResponseWriter, r *http.Request) {
	var articlesPublisherConfig ArticlesPublisherConfig

	if err := json.NewDecoder(r.Body).Decode(&articlesPublisherConfig); err != nil {
		errMsg := fmt.Errorf("invalid request: %w", err).Error()
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	pruningRepo := repository.NewPruning(toDomainAgent(articlesPublisherConfig.PruningAgent))
	articleRepo := repository.NewArticle(toDomainAgent(articlesPublisherConfig.ArticleAgent))
	publisherRepo := repository.NewPublisher()
	service := domain.NewServicePruning(pruningRepo, articleRepo, publisherRepo)

	err := service.PublishArticlesKeywords(toDomainCMS(articlesPublisherConfig.CMS), articlesPublisherConfig.Keywords, articlesPublisherConfig.ArticlePrompt, articlesPublisherConfig.PruningPrompt)
	if err != nil {
		errMsg := fmt.Errorf("error publishing articles: %w", err).Error()
		http.Error(w, errMsg, http.StatusInternalServerError)
	}
}
