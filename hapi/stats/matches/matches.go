package matches

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"go-halo.com/hapi/payloads"
)

// List returns some matches
func List() (payloads.MatchList, error) {
	url := os.Getenv("HAPI_URL")

	// The first step would be identifying a gamertag (gt)
	gt := "Supertigerman"
	mlp, err := json.Marshal(payloads.MatchListPayload{
		Gamertag: gt,
		Count:    1,
		Offset:   0,
		Mode:     "matchmade"})
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(mlp))
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	ml := payloads.MatchList{}
	json.Unmarshal([]byte(body), &ml)

	if true == false {
		return ml, errors.New("fake ass error")
	}

	return ml, nil
}
