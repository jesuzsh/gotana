package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	repo "github.com/ccamac01/halo-gofinite/halo-gofinite-service/repo"
)

const STATS_MATCHES_RETRIEVES_ENDPOINT = "https://dev--TestHalo.ccamacho.autocode.gg/"

func main() {
	resp, err := http.Get(STATS_MATCHES_RETRIEVES_ENDPOINT)
	if err != nil {
		log.Fatalf("failed to fetch match data")
	}

	defer resp.Body.Close()

	var matchDetails repo.InfiniteMatchDetailsResult
	err = json.NewDecoder(resp.Body).Decode(&matchDetails)
	if err != nil {
		log.Fatal(err)
	}
	prettyDetails, err := json.MarshalIndent(matchDetails, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(prettyDetails))
}
