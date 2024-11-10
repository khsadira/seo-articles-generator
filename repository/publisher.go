package repository

import (
	"fmt"

	"github.com/qantai/domain"
)

type Publisher struct{}

func NewPublisher() Publisher {
	return Publisher{}
}

func (p Publisher) PublishArticle(cms domain.CMS, article domain.Article) error {
	switch cms.ID {
	case "wordpress":
		article.Print()
		return nil
	default:
		return fmt.Errorf("cms not supported")
	}
}
