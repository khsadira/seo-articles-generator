package domain

type PublisherRepository interface {
	PublishArticle(cms CMS, article Article) error
}

type ArticleRepository interface {
	GenerateArticle(keyword, articlePrompt string) (Article, error)
}

type PruningRepository interface {
	GetPrunedKeywords(keyword, pruningPromt string) ([]string, error)
}
