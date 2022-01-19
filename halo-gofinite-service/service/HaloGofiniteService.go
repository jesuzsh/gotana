package service

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"go.uber.org/zap"

	"github.com/ccamac01/halo-gofinite/halo-gofinite-service/repo"
)

type HaloGofiniteService struct {
	statsMatchesEndpoint string
	log                  *zap.Logger
}

func InitLogger() *zap.Logger {
	logger, _ := zap.NewDevelopment()
	return logger
}

func NewHaloGofiniteService(statsMatchesEndpoint string) *HaloGofiniteService {
	log := InitLogger()

	return &HaloGofiniteService{
		statsMatchesEndpoint: statsMatchesEndpoint,
		log:                  log,
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

func (svc *HaloGofiniteService) GetMatchList() (string, error) {
	url := os.Getenv("HAPI_URL")

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("unable to obtain MatchList", err)
	}
	defer resp.Body.Close()

	var matchList repo.InfiniteMatchListResult
	err = json.NewDecoder(resp.Body).Decode(&matchList)
	if err != nil {
		return "", err
	}

	prettyList, err := json.MarshalIndent(matchList, "", "\t")
	if err != nil {
		return "", err
	}

	return string(prettyDetails), nil
}
