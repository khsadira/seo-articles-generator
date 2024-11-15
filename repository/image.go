package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/qantai/domain"
)

type ImageRepository struct {
	agent domain.ImageAgent
}

func NewImage(agent domain.ImageAgent) ImageRepository {
	return ImageRepository{
		agent: agent,
	}
}

func (i ImageRepository) GenerateImage(keyword, imagePrompt string) (domain.Image, error) {
	image, err := getImage(keyword, imagePrompt, i.agent)
	if err != nil {
		return domain.Image{}, fmt.Errorf("erreur lors de la récupération de l'image : %v", err)
	}

	return image, nil
}

func getImage(keyword string, imagePrompt string, imageAgent domain.ImageAgent) (domain.Image, error) {
	switch imageAgent.ID {
	case "openAI":
		url, err := getImageFromOpenAI(keyword, imagePrompt, imageAgent)
		if err != nil {
			return domain.Image{}, fmt.Errorf("erreur lors de la récupération de l'image : %v", err)
		}
		return domain.Image{
			ID:  keyword,
			URL: url,
		}, nil
	default:
		return domain.Image{}, nil
	}
}

func getImageFromOpenAI(keyword, prompt string, agent domain.ImageAgent) (string, error) {
	requestBody := map[string]interface{}{
		"prompt":  prompt + keyword,
		"model":   agent.Model,
		"n":       agent.N,
		"size":    agent.Size,
		"quality": agent.Quality,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/images/generations", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+agent.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute HTTP request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: %s", body)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	data, ok := response["data"].([]interface{})
	if !ok || len(data) == 0 {
		return "", fmt.Errorf("invalid response data format")
	}

	imageURL, ok := data[0].(map[string]interface{})["url"].(string)
	if !ok {
		return "", fmt.Errorf("failed to extract image URL")
	}

	return imageURL, nil
}
