package repo

import (
	"encoding/json"
	"log"
)

type MatchListPayload struct {
	Gamertag string `json:"gamertag"`
	Count    int    `json:"count"`
	Offset   int    `json:"offset"`
	Mode     string `json:"mode"`
}

func (mp *MatchListPayload) Marshal() []byte {
	payload, err := json.Marshal(mp)
	if err != nil {
		log.Printf("failed to marshal a MaxPayload", err)
	}

	return payload
}
