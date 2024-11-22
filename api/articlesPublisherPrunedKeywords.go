package api

import (
	"encoding/json"
	"net/http"

	"github.com/qantae/domain"
	"github.com/qantae/repository"
)

func HandlerPublishArticlesPrunedKeywords(w http.ResponseWriter, r *http.Request) {
	var articlesPublisherKeywordsConfig ArticlesPublisherPrunedKeywordsConfig

	if err := json.NewDecoder(r.Body).Decode(&articlesPublisherKeywordsConfig); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	articleRepo := repository.NewArticle(toDomainAgent(articlesPublisherKeywordsConfig.ArticleAgent))
	imageRepo := repository.NewImage(toDomainImageAgent(articlesPublisherKeywordsConfig.ImageAgent))
	publisherRepo := repository.NewPublisher()
	service := domain.NewService(articleRepo, imageRepo, publisherRepo)

	err := service.PublishArticlesPrunedKeywords(toDomainCMS(articlesPublisherKeywordsConfig.CMS), articlesPublisherKeywordsConfig.Keywords, articlesPublisherKeywordsConfig.ArticlePrompt, articlesPublisherKeywordsConfig.ImagePrompt)
	if err != nil {
		http.Error(w, "Error publishing articles", http.StatusInternalServerError)
	}
}
