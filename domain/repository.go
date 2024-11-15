package domain

type PublisherRepository interface {
	PublishArticle(cms CMS, article Article) error
	UploadImage(cms CMS, images Image) (UploadedImage, error)
}

type ArticleRepository interface {
	GenerateArticle(keyword, articlePrompt string, images []Image) (Article, error)
}

type PruningRepository interface {
	GetPrunedKeywords(keyword, pruningPromt string) ([]string, error)
}

type ImageRepository interface {
	GenerateImages(keyword, imagePrompt string, quantity int) ([]Image, error)
}
