package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
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

func matchListRequest(payload []byte) repo.InfiniteMatchListResult {
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

	ml := repo.InfiniteMatchListResult{}
	json.Unmarshal([]byte(body), &ml)

	return ml
}

// List returns some matches
func List(mlp repo.InfininteMatchListPayload) (repo.InfiniteMatchListResult, error) {
	if mlp.Gamertag == "" {
		return repo.InfiniteMatchListResult{}, errors.New("empty payload")
	}

	payload := mlp.Marshal()
	ml := matchListRequest(payload)

	return ml, nil
}

func Count(mlp repo.InfiniteMatchListPayload) (int, error) {
	if mlp.Gamertag == "" {
		return 0, errors.New("empty payload")
	}

	payload := mlp.Marshal()
	ml := matchListRequest(payload)

	return int(ml.Paging.Total), nil
}
