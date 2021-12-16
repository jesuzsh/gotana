package main

import (
	"fmt"
	"log"

	"go-halo.com/fetchers"
)

func main() {
	// The first step is deciding the gamertag, possibly hook this up to
	// command-line tooling.
	gamertag := "Lentilius"

	m := fetchers.Max()
	m.SetGamer(gamertag)

	ml, err := m.GetMatches()
	if err != nil {
		log.Fatalln(err)
	}

	for i, match := range ml.Data {
		fmt.Printf("%v %v:\n", i, match.Details.Map.Name)
	}

}
