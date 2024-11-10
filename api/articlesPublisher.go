package api

import (
	"encoding/json"
	"net/http"

	"github.com/qantai/domain"
)


type PublisherHandler struct {
	publisher domain.ServicePublisher
}

func NewPublisherHandler(publisher domain.ServicePublisher) *PublisherHandler {
	return &PublisherHandler{publisher: publisher}
}

func (h *PublisherHandler) PublishHandler(w http.ResponseWriter, r *http.Request) {
	var articlesPublisherConfig ArticlesPublisherConfig

	if err := json.NewDecoder(r.Body).Decode(&articlesPublisherConfig); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	for _, cms := range articlesPublisherConfig.CMS {
		if err := h.publisher.PublishArticlesPrunedKeywords(cms, articlesPublisherConfig.Keywords, articlesPublisherConfig.PrunedKeywords); err != nil {
			http.Error(w, "Failed to publish content", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Content published successfully"))
}
