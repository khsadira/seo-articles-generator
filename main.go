package main

import (
	"log"
	"net/http"

	"github.com/qantai/api"
	"github.com/qantai/domain"
	"github.com/qantai/repository"
)

func main() {
	handlePublishArticlesPrunedKeywords()
	handlePublishArticlesKeywords()

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handlePublishArticlesKeywords() {
	http.HandleFunc("/publishArticlesPrunedKeywords", api.HandlerPublishArticlesPrunedKeywords)
}

func handlePublishArticlesPrunedKeywords() {
	repo := repository.NewPublisher()
	service := domain.NewServicePublisher(repo)
	handler := api.NewPublisherHandler(service)

	http.HandleFunc("/publishArticlesKeywords", handler.PublishHandler)
}
