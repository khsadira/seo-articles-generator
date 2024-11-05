package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	numberOfCompetitorDomains  = "1"
	numberOfCompetitorKeywords = "1"
)

func getMostReleveantKeywords(keywords []string) ([]string, error) {
	var mostReleveantKeywordsResult []string

	for _, keyword := range keywords {
		println("keyword : ", keyword)
		baseURL := "https://api.semrush.com/"
		params := url.Values{}
		params.Add("type", "phrase_organic")
		params.Add("key", SEMRUSH_API_KEY)
		params.Add("phrase", keyword)
		params.Add("database", "fr")
		params.Add("display_limit", numberOfCompetitorDomains)

		apiURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
		resp, err := http.Get(apiURL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		links, err := getMostReleveantDomains(string(body))
		if err != nil {
			return nil, fmt.Errorf("aucune donnée trouvée pour le mot-clé : %s", keyword)
		}

		mostReleveantKeywords, err := getMostReleaventKeywords(links)
		if err != nil {
			return nil, err
		}

		mostReleveantKeywordsResult = append(mostReleveantKeywordsResult, mostReleveantKeywords...)
	}

	return removeDuplicates(mostReleveantKeywordsResult), nil
}

func getMostReleveantDomains(body string) ([]string, error) {
	reader := csv.NewReader(strings.NewReader(body))
	reader.Comma = ';'

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var domains []string

	for _, record := range records[1:] {
		println("most_relev:", record[0])
		domains = append(domains, record[0])
	}

	return domains, nil
}

func getMostReleaventKeywords(domains []string) ([]string, error) {
	var mostReleveantKeywordsResult []string

	for _, domain := range domains {
		println("domain : ", domain)
		baseURL := "https://api.semrush.com/"
		params := url.Values{}
		params.Add("type", "domain_organic")
		params.Add("key", SEMRUSH_API_KEY)
		params.Add("domain", domain)
		params.Add("database", "fr")
		params.Add("display_limit", numberOfCompetitorKeywords)

		apiURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
		resp, err := http.Get(apiURL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		keywords, err := getMostReleveantDomains(string(body))
		if err != nil {
			return nil, fmt.Errorf("aucune donnée trouvée pour le domaine : %s", domain)
		}

		mostReleveantKeywordsResult = append(mostReleveantKeywordsResult, keywords...)
	}

	return mostReleveantKeywordsResult, nil
}

func removeDuplicates(keywords []string) []string {
	seen := make(map[string]struct{})
	var result []string

	for _, word := range keywords {
		if _, exists := seen[word]; !exists {
			seen[word] = struct{}{}
			result = append(result, word)
		}
	}

	return result
}
