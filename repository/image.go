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
		URL: "https://cdn.discordapp.com/attachments/797558850076934144/1247677846898610256/image.png?ex=6735d04b&is=67347ecb&hm=fd3025b4fa6050cc3e4f0e0d65add9518adf6689bb47403fdfbebcb502c43e16&",
	}, nil
}
