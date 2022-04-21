package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jesuzsh/gotana/pkg/client"
	"github.com/jesuzsh/gotana/pkg/database"
)

const STATS_ENDPOINT = "STATS_ENDPOINT"

func intro() *database.User {
	var username string
	flag.StringVar(&username, "username", "", "gamertag")
	flag.Parse()

	if username == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Please enter a username: ")
		username, _ = reader.ReadString('\n')
		username = strings.TrimSuffix(username, "\n")

	}
	fmt.Printf("\n%v, welcome.\n", username)

	db := database.NewDevConnection()
	user, _ := db.CheckIn(username)

	return &user
}

func setup(user *database.User) *client.Client {
	statsEndpoint := os.Getenv(STATS_ENDPOINT)
	if statsEndpoint == "" {
		log.Fatal("please provide a valid autocode endpoint")
	}

	clt := client.NewClient(user, statsEndpoint)

	return clt
}

func main() {
	user := intro()
	clt := setup(user)

	fmt.Printf("\n Retrieving all match data...\n")
	clt.ProcessMatches()
	fmt.Printf("\n * Data processed. Check AWS.\n")

}
