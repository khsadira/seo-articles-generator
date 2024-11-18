package domain

import (
	"fmt"
	"log"
	"sync"
)

type Service struct {
	articleRepository   ArticleRepository
	imageRepository     ImageRepository
	publisherRepository PublisherRepository
}

type ServicePruning struct {
	pruningRepository   PruningRepository
	articleRepository   ArticleRepository
	imageRepository     ImageRepository
	publisherRepository PublisherRepository
}

func NewServicePruning(pruningRepository PruningRepository, articleRepository ArticleRepository, publisherRepository PublisherRepository) ServicePruning {
	return ServicePruning{
		pruningRepository:   pruningRepository,
		articleRepository:   articleRepository,
		publisherRepository: publisherRepository,
	}
}

func NewService(articleRepository ArticleRepository, imageRepository ImageRepository, publisherRepository PublisherRepository) Service {
	return Service{
		articleRepository:   articleRepository,
		imageRepository:     imageRepository,
		publisherRepository: publisherRepository,
	}
}

func (s ServicePruning) PublishArticlesKeywords(cms []CMS, keywords []string, articlePrompt, pruningPromt string) error {
	prunedKeywords, err := getPrunedKeywords(keywords, pruningPromt, s.pruningRepository)
	if err != nil {
		return fmt.Errorf("error while generating pruned keyword: %w", err)
	}

	articles, err := getArticles(cms, prunedKeywords, articlePrompt, "", s.articleRepository, s.imageRepository, s.publisherRepository)
	if err != nil {
		return fmt.Errorf("error while generating articles: %w", err)
	}

	publishArticles(s.publisherRepository, cms, articles)

	return nil
}

func (s Service) PublishArticlesPrunedKeywords(cms []CMS, keywords []string, articlePrompt, imagePrompt string) error {
	articles, err := getArticles(cms, keywords, articlePrompt, imagePrompt, s.articleRepository, s.imageRepository, s.publisherRepository)
	if err != nil {
		return fmt.Errorf("error while generating articles: %w", err)
	}

	// updateArticlesPlaceHolder(articles, imagePrompt, s.imageRepository)

	publishArticles(s.publisherRepository, cms, articles)

	return nil
}

func publishArticles(publisherRepository PublisherRepository, cms []CMS, articles []Article) {
	wg := sync.WaitGroup{}

	for _, cmsItem := range cms {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, article := range articles {
				err := publisherRepository.PublishArticle(cmsItem, article)
				if err != nil {
					log.Printf("Error publishing article to CMS: %v", err)
				}
			}
		}()
	}

	wg.Wait()
}
