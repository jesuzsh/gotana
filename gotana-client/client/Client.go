package service

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"go.uber.org/zap"

	"github.com/jesuzsh/gotana/gotana-client/repo"
)

type Client struct {
	StatsMatchesEndpoint   string
	StatsMatchListEndpoint string
	log                    *zap.Logger
	Buffer                 []uint8
}

func InitLogger() *zap.Logger {
	logger, _ := zap.NewDevelopment()
	return logger
}

func NewClient(statsMatchesEndpoint string, statsMatchListEndpoint string) *Client {
	log := InitLogger()

	return &Client{
		StatsMatchesEndpoint:   statsMatchesEndpoint,
		StatsMatchListEndpoint: statsMatchListEndpoint,
		log:                    log,
		Buffer:                 []uint8{},
	}
}

func (clt *Client) GetMatchDetails() (string, error) {
	resp, err := http.Get(clt.StatsMatchesEndpoint)
	log := clt.log
	if err != nil {
		log.Error("failed to fetch match data")
		return "", err
	}

	defer resp.Body.Close()

	var matchDetails repo.MatchDetailsResult
	err = json.NewDecoder(resp.Body).Decode(&matchDetails)
	if err != nil {
		return "", err
	}

	prettyDetails, err := json.MarshalIndent(matchDetails, "", "\t")
	if err != nil {
		return "", err
	}

	return string(prettyDetails), nil
}

func (clt *Client) GetMatchList(payload repo.MatchListPayload) (repo.MatchListResult, error) {
	log := clt.log
	if payload.Gamertag == "" {
		log.Fatal("No gamertag specified. Unable to get MatchList.")
	}

	resp, err := http.Post(
		clt.StatsMatchListEndpoint,
		"application/json",
		bytes.NewBuffer(payload.Marshal()),
	)
	if err != nil {
		log.Error("unable to obtain MatchList")
		return repo.MatchListResult{}, err
	}
	defer resp.Body.Close()

	var matchList repo.MatchListResult
	err = json.NewDecoder(resp.Body).Decode(&matchList)
	if err != nil {
		return repo.MatchListResult{}, err
	}

	return matchList, nil
}

func (clt *Client) TotalMatches() (int64, error) {
	payload := repo.MatchListPayload{
		Gamertag: "Lentilius",
		Count:    1,
		Offset:   0,
		Mode:     "matchmade",
	}

	mlr, err := clt.GetMatchList(payload)
	if err != nil {
		return 0, nil
	}

	return mlr.Paging.Total, nil
}

func (clt *Client) GetAllMatchList() {
	//pendingMatches, _ := clt.TotalMatches()
	pendingMatches := 50

	var wg sync.WaitGroup
	responses := make(chan repo.MatchListResult, pendingMatches)

	payload := repo.MatchListPayload{
		Gamertag: "Lentilius",
		Count:    25,
		Offset:   0,
		Mode:     "matchmade",
	}

	for pendingMatches > 0 {
		wg.Add(1)
		payload.Count = 25
		go func() {
			defer wg.Done()
			mlr, err := clt.GetMatchList(payload)
			if err != nil {
				log.Print(err)
				return
			}

			responses <- mlr
		}()

		pendingMatches -= 25
		payload.Offset += 25
	}

	go func() {
		wg.Wait()
		close(responses)
	}()

	for mlr := range responses {
		mlr.ListMatches()
	}

	return
}

func (clt *Client) WriteMatchList(filename string) (bool, error) {
	log := clt.log
	if len(clt.Buffer) == 0 {
		log.Fatal("Empty buffer. Nothing to write.")
	}

	err := ioutil.WriteFile(filename, clt.Buffer, 0644)
	if err != nil {
		log.Error("unable to save json file")
		return false, err
	}

	return true, nil
}
