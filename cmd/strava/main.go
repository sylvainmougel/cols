package main

import (
	"encoding/json"
	"fmt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"io"
	"net/http"
	"os"
)

const (
	stravaURL = "https://www.strava.com/api/v3"
	apiKey    = "d59bf111bc33f3f3e1cba117f7ed5d9c0e6db962"
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
	Data []float32 `json:"data"`
}

type Profile map[string]Series

type Segment struct {
	Name       string    `json:"name"`
	ObjectID   string    `json:"objectID"`
	SLAltitute []float32 `json:"slaltitude"`
	Distance   []float32 `json:"distance"`
	Slope      []float32 `json:"slope"`
	Tags       []string  `json:"_tags"`
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

func indexSegment(sti SegmentToIndex) {
	fmt.Println("Indexing: ", sti.Name)
	profile := getSegmentProfile(sti.SegmentID)
	infos := getSegmentInfo(sti.SegmentID)

	// Compute profile from see level
	distance := []float32{0.0}
	slAltitude := []float32{0.0}
	slope := []float32{0.0}
	altitudes := (*profile)["altitude"].Data
	distances := (*profile)["distance"].Data
	initialAlt := altitudes[0]
	lastDistance := float32(0.0)
	lastIndex := 0
	for i, d := range distances {
		if i == 0 {
			continue
		}
		// take only one point every 200m
		if d-lastDistance < 200 {
			continue
		}
		lastDistance = d
		slAltitude = append(slAltitude, altitudes[i]-initialAlt)
		distance = append(distance, d)
		da := altitudes[i] - altitudes[lastIndex]
		dd := distances[i] - distances[lastIndex]
		s := (da / dd) * 100
		slope = append(slope, s)
		lastIndex = i

	}

	segment := Segment{
		Name:       infos.Name,
		ObjectID:   sti.SegmentID,
		SLAltitute: slAltitude,
		Distance:   distance,
		Slope:      slope,
		Tags:       sti.Tags,
	}

	// Save in algolia
	client := search.NewClient("1QMZVCS1V5", "")
	// Create a new index and add a record
	index := client.InitIndex("profiles")
	resSave, err := index.SaveObjects(segment)
	if err != nil {
		fmt.Println(segment)
		panic(err)
	}
	resSave.Wait()

}

type SegmentToIndex struct {
	Name      string   `json:"name"`
	SegmentID string   `json:"objectID"`
	Tags      []string `json:"_tags"`
}

func main() {

	content, err := os.ReadFile("toIndex.json")
	if err != nil {
		panic(err)
	}
	var sti []SegmentToIndex
	err = json.Unmarshal(content, &sti)
	if err != nil {
		panic(err)
	}
	for _, s := range sti {
		indexSegment(s)
	}
}
