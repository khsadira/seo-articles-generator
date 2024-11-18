package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

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

func (i ImageRepository) GenerateImages(keyword, imagePrompt string, quantity int) ([]domain.Image, error) {
	images := make([]domain.Image, 0, quantity)

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	i.agent.N = 1
	for j := 0; j < quantity; j++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			generatedImages, err := getImages(keyword, imagePrompt, i.agent)
			if err != nil || len(generatedImages) == 0 {
				log.Printf("erreur lors de la récupération de l'image : %v", err)
				return
			}

			mu.Lock()
			images = append(images, generatedImages...)
			mu.Unlock()
		}()
	}

	wg.Wait()

	return images, nil
}

func getImages(keyword string, imagePrompt string, imageAgent domain.ImageAgent) ([]domain.Image, error) {
	switch imageAgent.ID {
	case "openAI":
		images, err := getImageFromOpenAI(keyword, imagePrompt, imageAgent)
		if err != nil {
			return nil, fmt.Errorf("erreur lors de la récupération de l'image : %v", err)
		}
		return images, nil
	case "mock":
		return getMockImage(keyword)
	default:
		return nil, fmt.Errorf("agent d'image non reconnu : %s", imageAgent.ID)
	}
}

func getMockImage(keyword string) ([]domain.Image, error) {
	return []domain.Image{{
		ID:  keyword,
		URL: "https://fastly.picsum.photos/id/736/200/300.jpg?hmac=WlU1DEqIVU_kIsTa682WsLgBIfCRbqhOAuKifGAq8TY",
	}}, nil
}

func getImageFromOpenAI(keyword, prompt string, agent domain.ImageAgent) ([]domain.Image, error) {
	requestBody := map[string]interface{}{
		"prompt":  prompt + keyword,
		"model":   agent.Model,
		"n":       agent.N,
		"size":    agent.Size,
		"quality": agent.Quality,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/images/generations", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+agent.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", body)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	data, ok := response["data"].([]interface{})
	if !ok || len(data) == 0 {
		return nil, fmt.Errorf("invalid response data format")
	}

	images := make([]domain.Image, len(data))

	for i, item := range data {
		imageURL, ok := item.(map[string]interface{})["url"].(string)
		if !ok {
			log.Printf("failed to extract image URL")
			continue
		}

		images[i] = domain.Image{
			URL: imageURL,
			ID:  keyword,
		}

	}

	return images, nil
}
