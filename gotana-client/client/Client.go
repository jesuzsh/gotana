package service

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
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

func (clt *Client) TotalMatches() (int, error) {
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

	return int(mlr.Paging.Total), nil
}

func (clt *Client) GetAllMatchList(out chan<- repo.MatchListResult, pendingMatches int) {
	// TODO: remove
	//pendingMatches = 25
	var wg sync.WaitGroup
	payload := repo.MatchListPayload{
		Gamertag: "Lentilius",
		Count:    25,
		Offset:   0,
		Mode:     "matchmade",
	}

	for pendingMatches > 0 {
		wg.Add(1)
		go func(p repo.MatchListPayload) {
			defer wg.Done()
			mlr, err := clt.GetMatchList(p)
			if err != nil {
				fmt.Println("there is an error")
				log.Print(err)
				return
			}

			out <- mlr
		}(payload)

		pendingMatches -= 25
		payload.Offset += 25

		if pendingMatches < 25 {
			payload.Count = pendingMatches
		}

	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return
}

func (clt *Client) ZipResults(out chan<- repo.ZipPayload, in <-chan repo.MatchListResult) {
	var wg sync.WaitGroup

	for mlr := range in {
		wg.Add(1)
		go func(result repo.MatchListResult) {
			defer wg.Done()
			result.ListMatches()

			fileJSON, err := json.Marshal(result)
			if err != nil {
				log.Fatal("json.Marshal:", err)
			}

			var fileGZ bytes.Buffer
			zipper := gzip.NewWriter(&fileGZ)

			_, err = zipper.Write(fileJSON)
			if err != nil {
				log.Fatalf("zipper.Write ERRR: %+v", err)
			}
			zipper.Close()

			zp := repo.ZipPayload{
				ID:  result.Data.ID,
				Zip: fileGZ,
			}
			out <- zp
		}(mlr)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return
}

func (clt *Client) Persister(in <-chan repo.ZipPayload) {
	for z := range in {
		// TODO
		fmt.Println("==========doing something with a zipfile============")
		fmt.Println(z)
	}
}

func (clt *Client) ProcessMatches() {
	totalMatches, _ := clt.TotalMatches()

	results := make(chan repo.MatchListResult, totalMatches/25)
	zippedResults := make(chan repo.ZipPayload, totalMatches/25)

	fmt.Println("GetAllMatchList")
	go clt.GetAllMatchList(results, totalMatches)
	fmt.Println("ZipResults")
	go clt.ZipResults(zippedResults, results)
	clt.Persister(zippedResults)

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
