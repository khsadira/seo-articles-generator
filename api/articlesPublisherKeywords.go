package api

import (
	"encoding/json"
	"net/http"

	"github.com/qantai/domain"
	"github.com/qantai/repository"
)

func HandlerPublishArticlesPrunedKeywords(w http.ResponseWriter, r *http.Request) {
	var articlesPublisherKeywordsConfig ArticlesPublisherKeywordsConfig

	if err := json.NewDecoder(r.Body).Decode(&articlesPublisherKeywordsConfig); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	pruningRepo := repository.NewPruning(toDomainAgent(articlesPublisherKeywordsConfig.PruningAgent))
	articleRepo := repository.NewArticle(toDomainAgent(articlesPublisherKeywordsConfig.ArticleAgent))
	publisherRepo := repository.NewPublisher()
	service := domain.NewServicePruning(pruningRepo, articleRepo, publisherRepo)

	err := service.PublishArticlesKeywords(toDomainCMS(articlesPublisherKeywordsConfig.CMS), articlesPublisherKeywordsConfig.Keywords)
	if err != nil {
		http.Error(w, "Error publishing articles", http.StatusInternalServerError)
	}
}
