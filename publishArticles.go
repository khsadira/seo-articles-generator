package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Article struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  string `json:"status"`
}

var API_URL string
var TOKEN string

func publishArticleCMS(cmsIDs []string, article Article) {
	for _, cmsID := range cmsIDs {
		switch cmsID {
		case "wordpress":
			publishArticleWP(article)
		}
	}
}

func publishArticleWP(article Article) error {
	API_URL = os.Getenv("WP_API_URL")
	if API_URL == "" {
		API_URL = "https://qantae.io/wp-json/wp/v2/posts"
	}

	TOKEN = os.Getenv("WP_TOKEN")
	if TOKEN == "" {
		TOKEN = "a2hhbi5zYWRpcmFjNDJAZ21haWwuY29tOnJLSmsgdkVUUSBlUExhIGhLQjkgZ2J6WSBmZ01I"
	}

	jsonData, err := json.Marshal(article)
	if err != nil {
		return fmt.Errorf("erreur lors de l'encodage en JSON: %v", err)
	}

	req, err := http.NewRequest("POST", API_URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Basic "+TOKEN)
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
