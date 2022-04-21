// all the Halo data associated with particular users.
package client

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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/jesuzsh/gotana/pkg/database"
	"github.com/jesuzsh/gotana/pkg/repo"
)

// Client contains essential attributes for accessing relevant endpoints.
type Client struct {
	User          *database.User
	StatsEndpoint string
}

// NewClient creates an instance of Client.
func NewClient(user *database.User, statsEndpoint string) *Client {
	return &Client{
		User:          user,
		StatsEndpoint: statsEndpoint,
	}
}

// GetMatchDetails retrieves the data associated with a particular match.
func (clt *Client) GetMatchDetails() (string, error) {
	resp, err := http.Get(clt.StatsEndpoint)
	if err != nil {
		log.Fatal("failed to fetch match data")
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

// GetMatchList retrieves match data for a particular user. This completes a
// single request to the associated endpoint. Request parameters are specified
// in a repo.MatchListPayload.
func (clt *Client) GetMatchList(payload repo.MatchListPayload) (repo.MatchListResult, error) {
	if payload.Gamertag == "" {
		log.Fatal("No gamertag specified. Unable to get MatchList.")
	}

	resp, err := http.Post(
		clt.StatsEndpoint,
		"application/json",
		bytes.NewBuffer(payload.Marshal()),
	)
	if err != nil {
		log.Fatal("unable to obtain MatchList")
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

// TotalMatches makes a single request to discover the total number of matches
// for a particular player.
func (clt *Client) TotalMatches() (int, error) {
	payload := repo.MatchListPayload{
		Gamertag: clt.User.Gamertag,
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

// GetAllMatchList retrieves all existing match data for a particular user.
// Many requests are made concurrently. The first stage of the pipeline.
func (clt *Client) GetAllMatchList(out chan<- repo.MatchListResult, pendingMatches int) {
	var wg sync.WaitGroup
	payload := repo.MatchListPayload{
		Gamertag: clt.User.Gamertag,
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

// ZipResults iterates over a channel that contains repo.MatchListResult's.
// This function performs the compression of json data. Compressed data is sent
// to an output channel for the next stage of the pipeline.
func (clt *Client) ZipResults(out chan<- repo.ZipPayload, in <-chan repo.MatchListResult) {
	var wg sync.WaitGroup

	for mlr := range in {
		wg.Add(1)
		go func(result repo.MatchListResult) {
			defer wg.Done()

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

// Persister is responsible for saving the json data in S3. The zipped files
// are uploaded to S3 using the aws-sdk-go.
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

			_, err = uploader.Upload(&s3manager.UploadInput{
				Bucket: aws.String("gotana"),
				Key:    aws.String(clt.User.Gamertag + "/" + zip.ID + ".json.gz"),
				Body:   file,
			})
			if err != nil {
				fmt.Printf("%+v", err)
				fmt.Errorf("failed to upload file, %v", err)
				return
			}
		}(z)
	}

	go func() {
		wg.Wait()
	}()

	return
}

// ProcessMatches is the overall orchestration of the pipeline. The required
// channels are created and sent to relevant stages of the pipeline.
func (clt *Client) ProcessMatches() {
	totalMatches, _ := clt.TotalMatches()

	results := make(chan repo.MatchListResult, totalMatches/25)
	zippedResults := make(chan repo.ZipPayload, totalMatches/25)

	go clt.GetAllMatchList(results, totalMatches)
	go clt.ZipResults(zippedResults, results)
	clt.Persister(zippedResults)

	db := database.NewDevConnection()
	db.MarkComplete(clt.User)

	return
}
