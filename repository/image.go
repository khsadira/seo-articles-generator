package repository

import (
	"fmt"

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
		return getImageFromOpenAI(keyword, imagePrompt, imageAgent)
	default:
		return domain.Image{}, nil
	}
}

func getImageFromOpenAI(keyword, _ string, _ domain.ImageAgent) (domain.Image, error) {
	return domain.Image{
		ID:  keyword,
		URL: "http://fakeimg.pl/250x100",
	}, nil
}
