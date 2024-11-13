package repository

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
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
		err := publishArticleWP(toArticle(article), cms.URL, cms.User, cms.APIKey)
		if err != nil {
			return fmt.Errorf("error while publishing article to wordpress: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("cms not supported")
	}
}

func publishArticleWP(article Article, url, user, pass string) error {
	jsonData, err := json.Marshal(article)
	if err != nil {
		return fmt.Errorf("erreur lors de l'encodage en JSON: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	apiKey := base64.StdEncoding.EncodeToString([]byte(user + ":" + pass))

	req.Header.Set("Authorization", "Basic "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("echec de la publication de l'article. Code d'erreur: %d | %s", resp.StatusCode, string(bodyBytes))
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
