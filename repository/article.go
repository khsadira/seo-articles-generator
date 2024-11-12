package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/qantai/domain"
)

type ArticleRepository struct {
	agent domain.Agent
}

func NewArticle(agent domain.Agent) ArticleRepository {
	return ArticleRepository{
		agent: agent,
	}
}

func (a ArticleRepository) GenerateArticle(keyword, articlePrompt string) (domain.Article, error) {
	title, content, err := getArticle(keyword, articlePrompt, a.agent)
	if err != nil {
		return domain.Article{}, fmt.Errorf("erreur lors de la récupération de l'article : %v", err)
	}

	return domain.Article{
		Title:   title,
		Content: content,
		Status:  "draft",
	}, nil
}

func getArticle(keyword string, articlePrompt string, agent domain.Agent) (string, string, error) {
	switch agent.ID {
	case "openAI":
		return getArticleFromOpenAI(keyword, articlePrompt, agent)
	default:
		return "", "", fmt.Errorf("agent non supporté : %s", agent.ID)
	}
}

type OpenAIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

type OpenAiResponseContent struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func getArticleFromOpenAI(keyword string, articlePrompt string, agent domain.Agent) (string, string, error) {
	prompt := fmt.Sprintf("renvoie moi sous format json {\"title\":\"{{titre_placeholder}}\",\"content\":\"{{content_placeHolder}}\"} la reponse pour le prompt suivant: %s", articlePrompt+keyword)

	requestBody := OpenAIRequest{
		Model:       agent.Model,
		Messages:    []Message{{Role: "user", Content: prompt}},
		Temperature: agent.Temperature,
		MaxTokens:   agent.MaxToken,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", "", err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Authorization", "Bearer "+agent.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("API error: %s", body)
	}

	var openAIResponse OpenAIResponse
	if err := json.Unmarshal(body, &openAIResponse); err != nil {
		return "", "", err
	}

	var contentResp OpenAiResponseContent
	if err := json.Unmarshal([]byte(openAIResponse.Choices[0].Message.Content), &contentResp); err != nil {
		return "", "", err
	}

	if contentResp.Title == "" || contentResp.Content == "" {
		return "", "", fmt.Errorf("invalid response format")
	}

	articleTitle := contentResp.Title
	articleContent := contentResp.Content

	return articleTitle, articleContent, nil
}
