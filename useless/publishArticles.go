package useless

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Article struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  string `json:"status"`
}

var (
	API_URL = ""
	TOKEN   = ""
)

func publishArticleCMS(cmsIDs []string, article Article) {
	for _, cmsID := range cmsIDs {
		switch cmsID {
		case "wordpress":
			publishArticleWP(article)
		}
	}
}

func publishArticleWP(article Article) error {
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
