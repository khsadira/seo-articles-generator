package domain

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"
)

type Article struct {
	Title         string
	Content       string
	Status        string
	FeaturedMedia float64
}

func (a Article) Print(cms CMS) {
	println(
		"________________________\n"+"CMS:", cms.ID+"_URL:"+cms.URL+"_APIKEY:"+cms.APIKey+"\n",
		"TITRE:", a.Title+"\n"+"CONTENT:", a.Content+"\n"+"STATUS:", a.Status+"\n"+"________________________\n",
	)
}

func getPrunedKeywords(keywords []string, pruningPromt string, pruningRepository PruningRepository) ([]string, error) {
	prunedKeywords := make([]string, 0)
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}

	for _, keyword := range keywords {
		wg.Add(1)
		go func() {
			defer wg.Done()

			prunedKeywordsTmp, err := pruningRepository.GetPrunedKeywords(keyword, pruningPromt)
			if err != nil {
				log.Printf("Error pruning keyword: %s", err.Error())
				return
			}

			mu.Lock()
			prunedKeywords = append(prunedKeywords, prunedKeywordsTmp...)
			mu.Unlock()
		}()
	}

	wg.Wait()

	return prunedKeywords, nil
}

func getArticles(cms []CMS, keywords []string, articlePrompt, imagePrompt string, articleRepository ArticleRepository, imageRepository ImageRepository, publisherRepository PublisherRepository) ([]Article, error) {
	articles := make([]Article, len(keywords))

	wg := sync.WaitGroup{}

	for i, keyword := range keywords {
		wg.Add(1)

		go func() {
			defer wg.Done()

			images, err := imageRepository.GenerateImages(keyword, imagePrompt, 1)
			if err != nil {
				log.Printf("Error generating images: %s", err.Error())
				return
			}

			uploadedImages, err := uploadImages(cms, images, publisherRepository)
			if err != nil || len(uploadedImages) == 0 {
				log.Printf("Error uploading images: %v", err)
				return
			}

			article, err := articleRepository.GenerateArticle(keyword, articlePrompt, nil)
			if err != nil {
				log.Printf("Error generating article: %s", err.Error())
				return
			}

			articlesGenerated := updateArticlesPlaceHolder(cms, []Article{article}, imagePrompt, imageRepository, publisherRepository)

			if len(articlesGenerated) > 0 {
				article = articlesGenerated[0]
			}

			article.Status = "draft"
			article.FeaturedMedia = uploadedImages[0].FeaturedMedia

			articles[i] = article
		}()
	}

	wg.Wait()

	return articles, nil
}

func uploadImages(cms []CMS, images []Image, publisherRepository PublisherRepository) ([]UploadedImage, error) {
	uploadedImages := make([]UploadedImage, 0)

	for _, cmsItem := range cms {
		for _, image := range images {
			uploadedImage, err := publisherRepository.UploadImage(cmsItem, image)
			if err != nil {
				log.Printf("Error uploading images: %s", err.Error())
				continue
			}

			uploadedImages = append(uploadedImages, uploadedImage)
		}

	}

	return uploadedImages, nil
}

func updateArticlesPlaceHolder(cms []CMS, articles []Article, imagePrompt string, imageRepository ImageRepository, publisherRepository PublisherRepository) []Article {
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

				uploadedImages, err := uploadImages(cms, images, publisherRepository)
				if err != nil {
					fmt.Printf("error uploading image: %s", err.Error())
					continue
				}

				if len(uploadedImages) > 0 {
					generatedImages[0].URL = uploadedImages[0].URL
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
