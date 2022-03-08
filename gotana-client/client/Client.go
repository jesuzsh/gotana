package service

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"go.uber.org/zap"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/jesuzsh/gotana/gotana-client/repo"
)

type Client struct {
	Target                 string
	StatsMatchesEndpoint   string
	StatsMatchListEndpoint string
	log                    *zap.Logger
	Buffer                 []uint8
}

func InitLogger() *zap.Logger {
	logger, _ := zap.NewDevelopment()
	return logger
}

func NewClient(target string, statsMatchesEndpoint string, statsMatchListEndpoint string) *Client {
	log := InitLogger()

	return &Client{
		Target:                 target,
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
		Gamertag: clt.Target,
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
		Gamertag: clt.Target,
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
			//result.ListMatches()

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
				ID:  strconv.Itoa(result.Paging.Offset),
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
	var wg sync.WaitGroup
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-2")},
	))
	uploader := s3manager.NewUploader(sess)

	for z := range in {
		wg.Add(1)
		go func(zip repo.ZipPayload) {
			defer wg.Done()
			file, err := ioutil.TempFile("", "prefix")
			if err != nil {
				log.Fatal(err)
			}
			defer os.Remove(file.Name())

			err = ioutil.WriteFile(file.Name(), []byte(zip.Zip.String()), 0644)
			if err != nil {
				fmt.Printf("WriteFileGZ ERROR: %+v", err)
			}

			fmt.Println("uploadiing to s3")
			result, err := uploader.Upload(&s3manager.UploadInput{
				Bucket: aws.String("gotana"),
				Key:    aws.String(clt.Target + "/" + zip.ID + ".json.gz"),
				Body:   file,
			})
			if err != nil {
				fmt.Printf("%+v", err)
				fmt.Errorf("failed to upload file, %v", err)
				return
			}

			fmt.Printf("file uploaded to, %s\n", aws.String(result.Location))

		}(z)
	}

	go func() {
		wg.Wait()
	}()

	return
}

func (clt *Client) ProcessMatches() {
	totalMatches, _ := clt.TotalMatches()

	results := make(chan repo.MatchListResult, totalMatches/25)
	zippedResults := make(chan repo.ZipPayload, totalMatches/25)

	go clt.GetAllMatchList(results, totalMatches)
	go clt.ZipResults(zippedResults, results)
	clt.Persister(zippedResults)

	fmt.Println("We done done.")

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
