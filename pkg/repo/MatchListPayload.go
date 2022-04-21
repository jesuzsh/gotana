package repo

import (
	"encoding/json"
	"log"
)

// MatchListPayload is the request parameters for acquiring a MatchListResult.
type MatchListPayload struct {
	Gamertag string `json:"gamertag"`
	Count    int    `json:"count"`
	Offset   int    `json:"offset"`
	Mode     string `json:"mode"`
}

// Marshal is a light rapper over the built-in json.Marshal(). It might make
// sense to get rid of this method as its responsibility is minimal.
func (mp *MatchListPayload) Marshal() []byte {
	payload, err := json.Marshal(mp)
	if err != nil {
		log.Printf("failed to marshal a MaxPayload", err)
	}

	return payload
}
