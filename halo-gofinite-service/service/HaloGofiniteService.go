package service

import (
	"bytes"
	"encoding/json"
	"fmt"
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
			Gamertag: "",
			Count:    3,
			Offset:   0,
			Mode:     "matchmade",
		},
		log: log,
	}
}

func (svc *HaloGofiniteService) SetGamer(gamertag string) {
	svc.StatsMatchListPayload.Gamertag = gamertag
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

func (svc *HaloGofiniteService) GetMatchList() (string, error) {
	resp, err := http.Post(
		svc.StatsMatchListEndpoint,
		"application/json",
		bytes.NewBuffer(svc.StatsMatchListPayload.Marshal()),
	)
	if err != nil {
		log.Printf("unable to obtain MatchList", err)
	}
	defer resp.Body.Close()

	var matchList repo.InfiniteMatchListResult
	fmt.Println("About to decode match list.")
	err = json.NewDecoder(resp.Body).Decode(&matchList)
	if err != nil {
		return "", err
	}
	fmt.Println("Completed the decode")

	prettyList, err := json.MarshalIndent(matchList, "", "\t")
	if err != nil {
		return "", err
	}

	return string(prettyList), nil
}
