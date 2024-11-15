package repository

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"

	imagePkg "image"
	_ "image/gif"  // Pour supporter le décodage GIF
	_ "image/jpeg" // Pour supporter le décodage JPEG
	_ "image/png"  // Pour supporter le décodage PNG

	"github.com/qantai/domain"
)

type Article struct {
	Title         string  `json:"title"`
	Content       string  `json:"content"`
	Status        string  `json:"status"`
	FeaturedMedia float64 `json:"featured_media"`
}

type Publisher struct{}

func NewPublisher() Publisher {
	return Publisher{}
}

func (p Publisher) PublishArticle(cms domain.CMS, article domain.Article) error {
	switch cms.ID {
	case "wordpress":
		err := publishArticleWP(toArticle(article), cms.URL, cms.User, cms.APIKey)
		if err != nil {
			return fmt.Errorf("error while publishing article to wordpress: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("cms not supported")
	}
}

func publishArticleWP(article Article, url, user, pass string) error {
	jsonData, err := json.Marshal(article)
	if err != nil {
		return fmt.Errorf("erreur lors de l'encodage en JSON: %v", err)
	}

	req, err := http.NewRequest("POST", url+"/posts", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	apiKey := base64.StdEncoding.EncodeToString([]byte(user + ":" + pass))

	req.Header.Set("Authorization", "Basic "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("echec de la publication de l'article. Code d'erreur: %d | %s", resp.StatusCode, string(bodyBytes))
	}

	fmt.Println("Article publié avec succès!")
	return nil
}

func toArticle(article domain.Article) Article {
	return Article{
		Title:         article.Title,
		Content:       article.Content,
		Status:        article.Status,
		FeaturedMedia: article.FeaturedMedia,
	}
}

func (p Publisher) UploadImage(cms domain.CMS, image domain.Image) (domain.UploadedImage, error) {

	switch cms.ID {
	case "wordpress":
		return uploadImageWP(image, cms.URL, cms.User, cms.APIKey)
	default:
		return domain.UploadedImage{}, fmt.Errorf("cms not supported")
	}
}

func convertImageAsBody(image domain.Image) (*bytes.Buffer, string, error) {
	resp, err := http.Get(image.URL)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("failed to get image, status code: %d", resp.StatusCode)
	}

	// Décoder l'image à partir de la réponse HTTP
	img, _, err := imagePkg.Decode(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode image: %w", err)
	}

	// Stocker l'image en PNG dans un buffer
	imageData := &bytes.Buffer{}
	err = png.Encode(imageData, img)
	if err != nil {
		return nil, "", fmt.Errorf("failed to encode image as PNG: %w", err)
	}

	// Créer le multipart/form-data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	imageName := filepath.Base(image.URL) + ".png" // Ajoute l'extension PNG

	part, err := writer.CreateFormFile("file", imageName)
	if err != nil {
		return nil, "", fmt.Errorf("error creating form file: %w", err)
	}

	_, err = part.Write(imageData.Bytes())
	if err != nil {
		return nil, "", fmt.Errorf("error writing image data to form: %w", err)
	}
	writer.Close()

	return body, writer.FormDataContentType(), nil
}

func uploadImageWP(image domain.Image, wpBaseURL, user, pass string) (domain.UploadedImage, error) {
	body, contentType, err := convertImageAsBody(image)
	if err != nil {
		return domain.UploadedImage{}, fmt.Errorf("error converting image to body: %w", err)
	}

	req, err := http.NewRequest("POST", wpBaseURL+"/media", body)
	if err != nil {
		return domain.UploadedImage{}, fmt.Errorf("error creating HTTP request: %w", err)
	}

	req.SetBasicAuth(user, pass)
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{}
	uploadResp, err := client.Do(req)
	if err != nil {
		return domain.UploadedImage{}, fmt.Errorf("error uploading image: %w", err)
	}
	defer uploadResp.Body.Close()

	if uploadResp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(uploadResp.Body)
		return domain.UploadedImage{}, fmt.Errorf("error uploading image status not created: %s", string(respBody))
	}

	var respData map[string]any

	if err := json.NewDecoder(uploadResp.Body).Decode(&respData); err != nil {
		return domain.UploadedImage{}, fmt.Errorf("error decoding response: %w", err)
	}

	id, ok := respData["id"].(float64)
	if !ok {
		return domain.UploadedImage{}, fmt.Errorf("error getting image id: %v", respData["id"])
	}

	return domain.UploadedImage{
		FeaturedMedia: id,
		URL:           respData["source_url"].(string),
	}, nil
}
