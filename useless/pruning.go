package useless

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

type KeywordData struct {
	Keyword      string
	SearchVolume int
	CPC          float64
	Competition  float64
	NbOfResult   int
	Intent       int
	Difficulty   int
}

func pruneCompetitorsKeywords(keywords []string) ([]string, error) {
	keywordsData, err := getKeywordsData(keywords)
	if err != nil {
		return nil, fmt.Errorf("aucune donnée trouvée pour les mot-clé : %s - %v", keywords, err)
	}

	scoredKeywords := make([]string, len(keywordsData))

	sort.Slice(keywordsData, func(i, j int) bool {
		return keywordsData[i].SearchVolume > keywordsData[j].SearchVolume
	})

	for _, data := range keywordsData {
		scoredKeywords = append(scoredKeywords, fmt.Sprintf("KEY: %s:\nSEARCH VOLUME: %d - CPC: %f - COMPETITION: %f - NUMBER OF RESULTS: %d - INTENT: %d - DIFFICULTY: %d\n\n",
			data.Keyword,
			data.SearchVolume,
			data.CPC,
			data.Competition,
			data.NbOfResult,
			data.Intent,
			data.Difficulty,
		))
	}

	return scoredKeywords, nil
}

const SEMRUSH_API_KEY = ""

func getKeywordsData(keywords []string) ([]KeywordData, error) {
	keywordsString := strings.Join(keywords, ";")

	baseURL := "https://api.semrush.com/"
	params := url.Values{}
	params.Add("type", "phrase_this")
	params.Add("key", SEMRUSH_API_KEY)
	params.Add("phrase", keywordsString)
	params.Add("database", "fr")
	params.Add("export_columns", "Dt,Db,Ph,Nq,Cp,Co,Nr,In,Kd") // A voir si on a besoin de tout ça

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

	keywordsData, err := parseKeywordsData(string(body))
	if err != nil {
		return nil, fmt.Errorf("aucune donnée trouvée pour les mot-clé : %s", keywordsString)
	}

	return keywordsData, nil
}

func parseKeywordsData(body string) ([]KeywordData, error) {
	reader := csv.NewReader(strings.NewReader(body))
	reader.Comma = ';'

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var data []KeywordData
	for _, record := range records[1:] {
		searchVolume, err := strconv.Atoi(record[1])
		if err != nil {
			return nil, fmt.Errorf("erreur en convertissant Search Volume pour %s : %v", record[0], err)
		}

		cpc, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, fmt.Errorf("erreur en convertissant CPC pour %s : %v", record[0], err)
		}

		competition, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			return nil, fmt.Errorf("erreur en convertissant Competition pour %s : %v", record[0], err)
		}

		nbOfResult, err := strconv.Atoi(record[4])
		if err != nil {
			return nil, fmt.Errorf("erreur en convertissant Number of Results pour %s : %v", record[0], err)
		}

		intent, err := strconv.Atoi(record[5])
		if err != nil {
			return nil, fmt.Errorf("erreur en convertissant Intent pour %s : %v", record[0], err)
		}

		difficulty, err := strconv.Atoi(record[6])
		if err != nil {
			return nil, fmt.Errorf("erreur en convertissant Difficulty pour %s : %v", record[0], err)
		}

		data = append(data, KeywordData{
			Keyword:      record[0],
			SearchVolume: searchVolume,
			CPC:          cpc,
			Competition:  competition,
			NbOfResult:   nbOfResult,
			Intent:       intent,
			Difficulty:   difficulty,
		})
	}

	return data, nil
}
