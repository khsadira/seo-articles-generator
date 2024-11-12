package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/qantai/domain"
)

type Article struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  string `json:"status"`
}

type Publisher struct{}

func NewPublisher() Publisher {
	return Publisher{}
}

func (p Publisher) PublishArticle(cms domain.CMS, article domain.Article) error {
	switch cms.ID {
	case "wordpress":
		publishArticleWP(toArticle(article), cms.URL, cms.APIKey)
		return nil
	default:
		return fmt.Errorf("cms not supported")
	}
}

func publishArticleWP(article Article, url, apiKey string) error {
	jsonData, err := json.Marshal(article)
	if err != nil {
		return fmt.Errorf("erreur lors de l'encodage en JSON: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Basic "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("echec de la publication de l'article. Code d'erreur: %d", resp.StatusCode)
	}

	fmt.Println("Article publié avec succès!")
	return nil
}

func toArticle(article domain.Article) Article {
	return Article{
		Title:   article.Title,
		Content: article.Content,
		Status:  "draft",
	}
}
