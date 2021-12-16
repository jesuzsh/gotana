package matches

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	pl "go-halo.com/hapi/payloads"
)

// List returns some matches
func List(mp pl.MaxPayload) (pl.MatchList, error) {
	if mp.Gamertag == "" {
		return pl.MatchList{}, errors.New("empty payload")
	}

	mlp, err := json.Marshal(mp)
	if err != nil {
		log.Printf("failed to marshal a MaxPayload", err)
	}

	url := os.Getenv("HAPI_URL")
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(mlp))
	if err != nil {
		log.Printf("unable to obtain MatchList", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("unable to read the response body", err)
	}

	ml := pl.MatchList{}
	json.Unmarshal([]byte(body), &ml)

	return ml, nil
}
