package domain

import (
	"fmt"
	"log"
	"regexp"
	"strings"
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

	articles, err := getArticles(cms, prunedKeywords, articlePrompt, s.articleRepository, s.imageRepository, s.publisherRepository)
	if err != nil {
		return fmt.Errorf("error while generating articles: %w", err)
	}

	publishArticles(s.publisherRepository, cms, articles)

	return nil
}

func (s Service) PublishArticlesPrunedKeywords(cms []CMS, keywords []string, articlePrompt, imagePromt string) error {
	articles, err := getArticles(cms, keywords, articlePrompt, s.articleRepository, s.imageRepository, s.publisherRepository)
	if err != nil {
		return fmt.Errorf("error while generating articles: %w", err)
	}

	updateArticlesPlaceHolder(articles, imagePromt, s.imageRepository)

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

func updateArticlesPlaceHolder(articles []Article, imagePrompt string, imageRepository ImageRepository) []Article {
	images := make([]Image, 0)

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	for _, article := range articles {
		wg.Add(1)
		go func() {
			defer wg.Done()

			placeHolders := findPlaceHolders(article.Content)

			for _, placeHolder := range placeHolders {
				if checkExisting(images, placeHolder.ID) || placeHolder.Type != PlaceHolderTypeImage {
					continue
				}

				generatedImages, err := imageRepository.GenerateImages(placeHolder.ID, imagePrompt, 1)
				if err != nil {
					fmt.Printf("error generating image: %s", err.Error())
					continue
				}

				mu.Lock()
				images = append(images, generatedImages[0])
				mu.Unlock()
			}

		}()
	}

	wg.Wait()

	return insertPlaceHolders(articles, images)
}

func insertPlaceHolders(articles []Article, images []Image) []Article {
	for i, article := range articles {
		for _, image := range images {
			article.Content = strings.ReplaceAll(article.Content, "{{"+image.ID+"_imageUrlPlaceHolder}}", image.URL)
		}
		articles[i] = article
	}

	return articles
}

func findPlaceHolders(str string) []PlaceHolder {
	re := regexp.MustCompile(`{{(.*?)_(.*?)UrlPlaceHolder}}`)
	matches := re.FindAllStringSubmatch(str, -1)
	var placeholders []PlaceHolder

	for _, match := range matches {
		placeholders = append(placeholders, PlaceHolder{
			ID:   match[1],
			Type: getPlaceHolerType(match[2]),
		})
	}

	return placeholders
}

func getPlaceHolerType(placeHolderType string) int {
	switch placeHolderType {
	case "image":
		return PlaceHolderTypeImage
	default:
		return -1
	}
}

func checkExisting(images []Image, id string) bool {
	for _, image := range images {
		if image.ID == id {
			return true
		}
	}
	return false
}
