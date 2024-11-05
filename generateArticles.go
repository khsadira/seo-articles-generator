package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func getArticleFromAgent(keyword string, agent string) Article {
	title, content, err := getArticle(keyword, agent)
	if err != nil {
		fmt.Println("Erreur lors de la récupération de l'article :", err)
		return Article{}
	}

	return Article{
		Title:   title,
		Content: content,
		Status:  "draft",
	}
}

func getArticle(keyword string, agent string) (string, string, error) {
	switch agent {
	case "openAI":
		return getArticleFromOpenAI(keyword)
	default:
		return "", "", fmt.Errorf("agent non supporté : %s", agent)
	}
}

const OPEN_AI_API_KEY = ""

// Structure pour la requête à l'API OpenAI
type OpenAIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

// Structure pour les messages dans la requête
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Structure pour la réponse de l'API OpenAI
type OpenAIResponse struct {
	Choices []Choice `json:"choices"`
}

// Structure pour les choix dans la réponse
type Choice struct {
	Message Message `json:"message"`
}

func getArticleFromOpenAI(keyword string) (string, string, error) {
	prompt := fmt.Sprintf("genre un article SEO optimisé pour le keyword %s", keyword)

	// Préparer la requête
	requestBody := OpenAIRequest{
		Model:       "gpt-3.5-turbo",
		Messages:    []Message{{Role: "user", Content: prompt}},
		Temperature: 0.7,
		MaxTokens:   500,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", "", err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Authorization", "Bearer "+OPEN_AI_API_KEY)
	req.Header.Set("Content-Type", "application/json")

	// Envoyer la requête à l'API OpenAI
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	// Lire la réponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	// Vérifier le statut de la réponse
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("API error: %s", body)
	}

	// Décodez la réponse JSON
	var openAIResponse OpenAIResponse
	if err := json.Unmarshal(body, &openAIResponse); err != nil {
		return "", "", err
	}

	articleTitle := fmt.Sprintf("Article sur %s", keyword)
	articleContent := openAIResponse.Choices[0].Message.Content

	return articleTitle, articleContent, nil
}
