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

func matchListRequest(payload []byte) pl.MatchList {
	url := os.Getenv("HAPI_URL")

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
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

	return ml
}

// List returns some matches
func List(mp pl.MaxPayload) (pl.MatchList, error) {
	if mp.Gamertag == "" {
		return pl.MatchList{}, errors.New("empty payload")
	}

	payload := mp.Marshal()
	ml := matchListRequest(payload)

	return ml, nil
}

func Count(mp pl.MaxPayload) (int64, error) {
	if mp.Gamertag == "" {
		return 0, errors.New("empty payload")
	}

	payload := mp.Marshal()
	ml := matchListRequest(payload)

	return ml.Paging.Total, nil
}
