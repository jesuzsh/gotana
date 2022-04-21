package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jesuzsh/gotana/pkg/database"
)

const STATS_MATCH_LIST_ENDPOINT = "STATS_MATCH_LIST_ENDPOINT"

func intro() {
	var user string
	flag.StringVar(&user, "user", "", "gamertag")
	flag.Parse()

	if user == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Please enter a user: ")
		user, _ = reader.ReadString('\n')
		user = strings.TrimSuffix(user, "\n")

	}
	fmt.Printf("\n%v, welcome.\n", user)
}

func setup() {
	statsMatchListEndpoint := os.Getenv(STATS_MATCH_LIST_ENDPOINT)

	if statsMatchListEndpoint == "" {
		log.Fatal("please provide a valid autocode endpoint")
	}

	//clt := client.NewClient("Lentilius", "", statsMatchListEndpoint)

	//clt.ProcessMatches()

}

func main() {
	intro()

	database.NewDevConnection()
}
