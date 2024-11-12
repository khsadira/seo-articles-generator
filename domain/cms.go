package domain

type CMS struct {
	ID     string
	URL    string
	APIKey string
}

type ServicePublisher struct {
	repository PublisherRepository
}

func NewServicePublisher(repository PublisherRepository) ServicePublisher {
	return ServicePublisher{repository: repository}
}

func (s ServicePublisher) PublishArticles(cms CMS, articles []Article) error {
	for _, article := range articles {
		err := s.repository.PublishArticle(cms, article)
		if err != nil {
			return err
		}
	}

	return nil
}
