package main

import (
	"log"
	"net/http"

	"github.com/qantai/api"
)

func main() {
	handlePublishArticles()
	handlePublishArticlesPrunedKeywords()

	port := ":8080"

	log.Println("Starting server on " + port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func handlePublishArticles() {
	http.HandleFunc("/publishArticles", api.HandlerPublishArticles)
}

func handlePublishArticlesPrunedKeywords() {
	http.HandleFunc("/publishArticlesPrunedKeywords", api.HandlerPublishArticlesPrunedKeywords)
}
