package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"go.uber.org/zap"

	"github.com/ccamac01/halo-gofinite/halo-gofinite-service/repo"
)

type HaloGofiniteService struct {
	statsMatchesEndpoint   string
	StatsMatchListEndpoint string
	StatsMatchListPayload  *repo.InfiniteMatchListPayload
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
		StatsMatchListPayload: &repo.InfiniteMatchListPayload{
			Gamertag: "Lentilius",
			Count:    25,
			Offset:   0,
			Mode:     "matchmade",
		},
		log:    log,
		Buffer: []uint8{},
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

func (svc *HaloGofiniteService) GetMatchList() (repo.InfiniteMatchListResult, error) {
	log := svc.log
	if svc.StatsMatchListPayload.Gamertag == "" {
		log.Fatal("No gamertag specified. Unable to get MatchList.")
	}

	resp, err := http.Post(
		svc.StatsMatchListEndpoint,
		"application/json",
		bytes.NewBuffer(svc.StatsMatchListPayload.Marshal()),
	)
	if err != nil {
		log.Error("unable to obtain MatchList")
		return repo.InfiniteMatchListResult{}, err
	}
	defer resp.Body.Close()

	var matchList repo.InfiniteMatchListResult
	err = json.NewDecoder(resp.Body).Decode(&matchList)
	if err != nil {
		return repo.InfiniteMatchListResult{}, err
	}

	// TODO temporary way to access the data
	svc.Buffer, err = json.MarshalIndent(matchList, "", "\t")
	if err != nil {
		return repo.InfiniteMatchListResult{}, err
	}

	return matchList, nil
}

func (svc *HaloGofiniteService) TotalMatches() (int64, error) {
	svc.StatsMatchListPayload.Count = 1

	mlr, _ := svc.GetMatchList()

	return mlr.Paging.Total, nil
}

// TODO: Implement this function in go routines.
func (svc *HaloGofiniteService) GetAllMatchList() ([]repo.InfiniteMatchListResult, error) {
	var mlrList []repo.InfinitMatchListResult
	pendingMatches, _ := svc.TotalMatches()

	for pendingMatches > 0 {
		fmt.Printf("Pending matches: %v\n", pendingMatches)
		svc.StatsMatchListPayload.Count = 25
		mlr, err := svc.GetMatchList()
		if err != nil {
			log.Fatal(err)
		}
		mlrList = append(mlrList, mlr)

		pendingMatches -= 25
		svc.StatsMatchListPayload.Offset += 25
	}

	return mlrList, nil
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
