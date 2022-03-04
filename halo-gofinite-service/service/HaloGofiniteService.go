package service

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"go.uber.org/zap"

	"github.com/ccamac01/halo-gofinite/halo-gofinite-service/repo"
)

type HaloGofiniteService struct {
	statsMatchesEndpoint   string
	StatsMatchListEndpoint string
	log                    *zap.Logger
	Buffer                 []uint8
}

func InitLogger() *zap.Logger {
	logger, _ := zap.NewDevelopment()
	return logger
}

func NewHaloGofiniteService(statsMatchesEndpoint string, statsMatchListEndpoint string) *HaloGofiniteService {
	log := InitLogger()

	return &HaloGofiniteService{
		statsMatchesEndpoint:   statsMatchesEndpoint,
		StatsMatchListEndpoint: statsMatchListEndpoint,
		log:                    log,
		Buffer:                 []uint8{},
	}
}

func (svc *HaloGofiniteService) GetMatchDetails() (string, error) {
	resp, err := http.Get(svc.statsMatchesEndpoint)
	log := svc.log
	if err != nil {
		log.Error("failed to fetch match data")
		return "", err
	}

	defer resp.Body.Close()

	var matchDetails repo.InfiniteMatchDetailsResult
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

func (svc *HaloGofiniteService) GetMatchList(payload repo.MatchListPayload) (repo.MatchListResult, error) {
	log := svc.log
	if payload.Gamertag == "" {
		log.Fatal("No gamertag specified. Unable to get MatchList.")
	}

	resp, err := http.Post(
		svc.StatsMatchListEndpoint,
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

func (svc *HaloGofiniteService) TotalMatches() (int64, error) {
	payload := repo.MatchListPayload{
		Gamertag: "Lentilius",
		Count:    1,
		Offset:   0,
		Mode:     "matchmade",
	}

	mlr, err := svc.GetMatchList(payload)
	if err != nil {
		return 0, nil
	}

	return mlr.Paging.Total, nil
}

func (svc *HaloGofiniteService) GetAllMatchList() {
	//pendingMatches, _ := svc.TotalMatches()
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
			mlr, err := svc.GetMatchList(payload)
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

func (svc *HaloGofiniteService) WriteMatchList(filename string) (bool, error) {
	log := svc.log
	if len(svc.Buffer) == 0 {
		log.Fatal("Empty buffer. Nothing to write.")
	}

	err := ioutil.WriteFile(filename, svc.Buffer, 0644)
	if err != nil {
		log.Error("unable to save json file")
		return false, err
	}

	return true, nil
}
