package domain

import (
	"log"
	"sync"
)

type Article struct {
	Title   string
	Content string
	Status  string
}

func (a Article) Print() {
	println(
		"________________________\n"+"TITRE:", a.Title+"\n"+"CONTENT:", a.Content+"\n"+"STATUS:", a.Status+"\n"+"________________________\n",
	)
}

func getPrunedKeywords(keywords []string, pruningRepository PruningRepository) ([]string, error) {
	prunedKeywords := make([]string, 0)
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}

	for _, keyword := range keywords {
		wg.Add(1)
		go func() {
			defer wg.Done()

			prunedKeywordsTmp, err := pruningRepository.GetPrunedKeywords(keyword)
			if err != nil {
				log.Printf("Error pruning keyword: %s", err.Error())
				return
			}

			mu.Lock()
			prunedKeywords = append(prunedKeywords, prunedKeywordsTmp...)
			mu.Unlock()
		}()
	}

	return prunedKeywords, nil
}

func getArticles(keywords []string, articleRepository ArticleRepository) ([]Article, error) {
	articles := make([]Article, len(keywords))

	wg := sync.WaitGroup{}

	for i, keyword := range keywords {
		wg.Add(1)

		go func() {
			defer wg.Done()

			article, err := articleRepository.GenerateArticle(keyword)
			if err != nil {
				log.Printf("Error generating article: %s", err.Error())
				return
			}

			articles[i] = article
		}()
	}

	wg.Wait()

	return articles, nil
}
