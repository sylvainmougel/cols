package main

import (
	"encoding/json"
	"fmt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"io"
	"net/http"
)

const (
	stravaURL = "https://www.strava.com/api/v3"
	apiKey    = "11565fd054b1192df14c0c0a2db8128dac27cda1"
)

func callStravaApi(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf(fmt.Sprintf("Bearer %s", apiKey)))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err

	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return body, nil
}

type Series struct {
	Data []float64 `json:"data"`
}

type Profile map[string]Series

type Segment struct {
	Name     string  `json:"name"`
	ObjectID string  `json:"objectID"`
	Profile  Profile `json:"profile"`
}

func getSegmentProfile(segmentID string) *Profile {
	resp, err := callStravaApi(fmt.Sprintf("%s/segments/%s/streams?keys=altitude&key_by_type=true", stravaURL, segmentID))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var p Profile
	err = json.Unmarshal(resp, &p)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &p
}

type Infos struct {
	Name string `json:"name"`
}

func getSegmentInfo(segmentID string) *Infos {
	// 'https://www.strava.com/api/v3/segments/17275870'
	resp, _ := callStravaApi(fmt.Sprintf("%s/segments/%s", stravaURL, segmentID))
	var i Infos
	err := json.Unmarshal(resp, &i)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &i
}

func main() {
	segmentID := "5211636"
	profile := getSegmentProfile(segmentID)
	infos := getSegmentInfo(segmentID)
	fmt.Println(profile)

	// Save in algolia
	client := search.NewClient("1QMZVCS1V5", "")

	segment := Segment{Profile: *profile, Name: infos.Name, ObjectID: string(segmentID)}
	// Create a new index and add a record
	index := client.InitIndex("profiles")
	resSave, err := index.SaveObjects(segment)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	resSave.Wait()
	//getSegmentInfo(17275870)

}
