package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	JSON_TYPE                 = "application/json"
	ACCEPT_HTTP_HEADER        = "accept"
	CONTENT_TYPE_HTTP_HEADER  = "content-type"
	AUTHORIZATION_HTTP_HEADER = "authorization"
)

func main() {
	url := "https://clientapi.pricingsaas.com/getExploreCompanies"
	token := os.Getenv("API_TOKEN")
	httpClient := http.DefaultClient

	for page := 1; page <= 122; page++ {
		payload := fmt.Sprintf("{\"page\":%d,\"filter\":{\"categories\":[],\"models\":[],\"pricingModels\":[],\"pricingMetrics\":[],\"sizes\":[]}}", page)

		req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(payload))
		if err != nil {
			log.Fatalf("error occurred while forming request: %v", err)
		}

		req.Header.Add(ACCEPT_HTTP_HEADER, JSON_TYPE)
		req.Header.Add(CONTENT_TYPE_HTTP_HEADER, JSON_TYPE)
		req.Header.Add(AUTHORIZATION_HTTP_HEADER, fmt.Sprintf("Bearer %s", token))

		res, err := httpClient.Do(req)
		if err != nil {
			log.Printf("error occurred while sending request: %v", err)
		}

		if res.StatusCode != 200 {
			defer res.Body.Close()
			body, _ := io.ReadAll(res.Body)
			log.Printf("error occurred: expected 201 created response but got: \nstatus code: %v, status: %v, body: %v", res.StatusCode, res.Status, string(body))
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Printf("error occurred while reading response body: %v", err)
		}

		filepath := fmt.Sprintf("page-%d.json", page)
		err = os.WriteFile(filepath, body, os.ModePerm)
		if err != nil {
			log.Printf("error occurred while writing to %s json file: %v", filepath, err)
		}

		err = res.Body.Close()
		if err != nil {
			log.Printf("error occurred while closing response body: %v", err)
		}

		log.Printf("page %d done!! :)", page)
	}

}
