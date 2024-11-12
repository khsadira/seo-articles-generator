package main

import (
	"log"
	"net/http"

	"github.com/qantai/api"
)

func main() {
	handlePublishArticles()
	handlePublishArticlesPrunedKeywords()

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handlePublishArticles() {
	http.HandleFunc("/publishArticles", api.HandlerPublishArticles)
}

func handlePublishArticlesPrunedKeywords() {
	http.HandleFunc("/publishArticlesPrunedKeywords", api.HandlerPublishArticlesPrunedKeywords)
}
