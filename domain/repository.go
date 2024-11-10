package domain

type PublisherRepository interface {
	PublishArticle(cms CMS, article Article) error
}

type ArticleRepository interface {
	GenerateArticle(keyword string) (Article, error)
}

type PruningRepository interface {
	GetPrunedKeywords(keyword string) ([]string, error)
}
