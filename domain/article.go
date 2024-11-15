package domain

import (
	"log"
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

func getArticles(cms []CMS, keywords []string, articlePrompt string, articleRepository ArticleRepository, imageRepository ImageRepository, publisherRepository PublisherRepository) ([]Article, error) {
	articles := make([]Article, len(keywords))

	wg := sync.WaitGroup{}

	for i, keyword := range keywords {
		wg.Add(1)

		go func() {
			defer wg.Done()

			images, err := imageRepository.GenerateImages(keyword, articlePrompt, 2)
			if err != nil {
				log.Printf("Error generating images: %s", err.Error())
				return
			}

			uploadedImages, err := uploadImages(cms, images, publisherRepository)
			if err != nil || len(uploadedImages) == 0 {
				log.Printf("Error uploading images: %s", err.Error())
				return
			}

			article, err := articleRepository.GenerateArticle(keyword, articlePrompt, nil)
			if err != nil {
				log.Printf("Error generating article: %s", err.Error())
				return
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
