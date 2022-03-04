package main

import (
	"fmt"
	"os"

	client "github.com/jesuzsh/gotana/gotana-client/client"
)

const STATS_MATCHES_RETRIEVES_ENDPOINT = "STATS_MATCHES_RETRIEVES_ENDPOINT"
const STATS_MATCH_LIST_ENDPOINT = "STATS_MATCH_LIST_ENDPOINT"

// TODO: add context to pass logger to other methods
func main() {
	log := client.InitLogger()
	statsMatchesEndpoint := os.Getenv(STATS_MATCHES_RETRIEVES_ENDPOINT)
	statsMatchListEndpoint := os.Getenv(STATS_MATCH_LIST_ENDPOINT)

	// TODO: add better input sanitization to avoid making unnecessary API calls
	if statsMatchesEndpoint == "" || statsMatchListEndpoint == "" {
		log.Fatal("please provide a valid autocode endpoint for retrieving match stats")
	}

	clt := client.NewHaloGofiniteService(statsMatchesEndpoint, statsMatchListEndpoint)

	clt.GetAllMatchList()
	fmt.Println("Done.")

}