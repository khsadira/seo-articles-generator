package domain

import (
	"fmt"
	"log"
)

type Agent struct {
	ID          string
	APIKey      string
	Temperature float64
}

type Service struct {
	articleRepository   ArticleRepository
	publisherRepository PublisherRepository
}

type ServicePruning struct {
	pruningRepository   PruningRepository
	articleRepository   ArticleRepository
	publisherRepository PublisherRepository
}

func NewServicePruning(pruningRepository PruningRepository, articleRepository ArticleRepository, publisherRepository PublisherRepository) ServicePruning {
	return ServicePruning{
		pruningRepository:   pruningRepository,
		articleRepository:   articleRepository,
		publisherRepository: publisherRepository,
	}
}

func NewService(articleRepository ArticleRepository, publisherRepository PublisherRepository) Service {
	return Service{
		articleRepository:   articleRepository,
		publisherRepository: publisherRepository,
	}
}

func (s ServicePruning) PublishArticlesKeywords(cms []CMS, keywords []string) error {
	prunedKeywords, err := getPrunedKeywords(keywords, s.pruningRepository)
	if err != nil {
		return fmt.Errorf("error while generating pruned keyword: %w", err)
	}

	articles, err := getArticles(prunedKeywords, s.articleRepository)
	if err != nil {
		return fmt.Errorf("error while generating articles: %w", err)
	}

	publishArticles(s.publisherRepository, cms, articles)

	return nil
}

func (s Service) PublishArticlesPrunedKeywords(cms []CMS, keywords []string) error {
	articles, err := getArticles(keywords, s.articleRepository)
	if err != nil {
		return fmt.Errorf("error while generating articles: %w", err)
	}

	publishArticles(s.publisherRepository, cms, articles)

	return nil
}

func publishArticles(publisherRepository PublisherRepository, cms []CMS, articles []Article) {
	for _, cmsItem := range cms {
		go func() {
			for _, article := range articles {
				err := publisherRepository.PublishArticle(cmsItem, article)
				if err != nil {
					log.Printf("Error publishing article to CMS: %v", err)
				}
			}
		}()
	}
}
