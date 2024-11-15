package domain

type Image struct {
	ID  string
	URL string
}

type UploadedImage struct {
	URL           string
	FeaturedMedia float64
}
