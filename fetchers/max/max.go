// Every gamertag will have a max.
package max

import (
	"log"

	pl "go-halo.com/hapi/payloads"
	"go-halo.com/hapi/stats/matches"
)

type Max struct {
	Payload pl.MaxPayload
	Matches []pl.MatchList
}

func (max *Max) SetGamer(gt string) {
	max.Payload.Gamertag = gt
	max.Payload.Count = 3
	max.Payload.Offset = 0
	max.Payload.Mode = "matchmade"
}

func (max *Max) GetMatches() (pl.MatchList, error) {
	ml, err := matches.List(max.Payload)
	if err != nil {
		log.Printf("Max could not fetch matches: %v", err)
		return ml, err
	}

	return ml, nil
}

func (max *Max) GetAllMatches() ([]pl.MatchList, error) {
	total, err := max.GetNumMatches()
	if err != nil {
		log.Printf("Max failed to get all the matches: %v", err)
		return []pl.MatchList{}, err
	}

	max.Payload.Count = 15
	max.Payload.Offset = 0

	foundMatches := make(chan pl.MatchList, total)
	for max.Payload.Offset < total {
		go func() {
			ml, err := max.GetMatches()
			if err != nil {
				log.Printf("error when attempting to get all matches: %v", err)
			}

			max.Matches = append(max.Matches, ml)

			//foundMatches <- ml

		}()

		max.Payload.Offset += max.Payload.Count

	}

	return max.Matches, nil

}

func collectMatches(foundMatches <-chan pl.MatchList) {

}

func (m *Max) GetNumMatches() (int, error) {
	count, err := matches.Count(m.Payload)
	if err != nil {
		log.Printf("Max couldn't count the number of matches: %v", err)
		return count, err
	}

	return count, nil
}
