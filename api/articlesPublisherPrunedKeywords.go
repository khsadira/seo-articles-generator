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

	articleRepo := repository.NewArticle(toDomainAgent(articlesPublisherKeywordsConfig.ArticleAgent))
	publisherRepo := repository.NewPublisher()
	service := domain.NewService(articleRepo, publisherRepo)

	err := service.PublishArticlesPrunedKeywords(toDomainCMS(articlesPublisherKeywordsConfig.CMS), articlesPublisherKeywordsConfig.Keywords, articlesPublisherKeywordsConfig.ArticlePrompt)
	if err != nil {
		http.Error(w, "Error publishing articles", http.StatusInternalServerError)
	}
}
