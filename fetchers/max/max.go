// Every gamertag will have a max.
package max

import (
	"log"

	pl "go-halo.com/hapi/payloads"
	"go-halo.com/hapi/stats/matches"
)

type Max struct {
	payload pl.MaxPayload
}

func (m *Max) SetGamer(gt string) {
	m.payload.Gamertag = gt
	m.payload.Count = 3
	m.payload.Offset = 0
	m.payload.Mode = "matchmade"
}

func (m *Max) GetMatches() (pl.MatchList, error) {
	ml, err := matches.List(m.payload)
	if err != nil {
		log.Printf("Max could not fetch matches: %v", err)
		return ml, err
	}

	return ml, nil
}
