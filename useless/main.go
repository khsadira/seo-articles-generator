package useless

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Print("starting server...")
	http.HandleFunc("/publishArticles", handlerPublishArticles)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
		log.Printf("defaulting to port %s", port)
	}

	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

type RequestBody struct {
	Keywords       []string `json:"keywords"`
	PrunedKeywords []string `json:"prunedKeywords"`
	CMS            []string `json:"cms"`
}

func handlerPublishArticles(w http.ResponseWriter, r *http.Request) {
	var body RequestBody

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	seoKeywords, err := getMostReleveantKeywords(body.Keywords)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	keywords := append(body.PrunedKeywords, seoKeywords...)

	agent := os.Getenv("agentAI")
	if agent == "" {
		agent = "openAI"
	}

	for _, keyword := range keywords {
		article := getArticleFromAgent(keyword, agent)

		publishArticleCMS(body.CMS, article)
	}
}

func handlerPublishArticles2(w http.ResponseWriter, r *http.Request) {
}
