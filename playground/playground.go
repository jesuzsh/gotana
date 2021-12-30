package main

import (
	"fmt"
	"log"

	"go-halo.com/fetchers"
)

func main() {
	// The first step is deciding the gamertag, possibly hook this up to
	// command-line tooling.
	gamertag := "Killamannjaro"

	//fetchMatches(gamertag)
	//fetchTotalMatches(gamertag)
	fetchAllMatches(gamertag)
}

func fetchAllMatches(gamertag string) {
	max := fetchers.Max()
	max.SetGamer(gamertag)

	matches, err := max.GetAllMatches()
	if err != nil {
		log.Printf("unable to get all matches: %v", err)
	}

	log.Printf("Length of matches: %v, Type: %T\n", len(matches), matches)
	log.Println(matches)
}

func fetchMatches(gamertag string) {
	max := fetchers.Max()
	max.SetGamer(gamertag)

	ml, err := max.GetMatches()
	if err != nil {
		log.Printf("unable to get matches: %v", err)
	}

	for i, match := range ml.Data {
		fmt.Printf("%v %v:\n", i, match.Details.Map.Name)
	}
}

func fetchTotalMatches(gamertag string) {
	max := fetchers.Max()
	max.SetGamer(gamertag)

	count, err := max.GetNumMatches()
	if err != nil {
		log.Printf("unable to get total matches: %v", err)
	}

	log.Printf("%v:\n%v total matches\n", gamertag, count)
}
