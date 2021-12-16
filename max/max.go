package main

import (
	"fmt"

	"go-halo.com/hapi/stats/matches"
)

func main() {
	ml, _ := matches.List()

	for i, match := range ml.Data {
		fmt.Printf("%v %v (Rank - %v):\n", i, match.Details.Map.Name, match.Player.Rank)
	}
}
