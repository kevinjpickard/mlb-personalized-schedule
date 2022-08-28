package schedule

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	DateFormat = "2006-01-02"
	sportID          = "1"
	scheduleEndpoint = "https://statsapi.mlb.com/api/v1/schedule"
)

type Schedule struct {
	TotalItems           int     `json:"totalItems"`
	TotalEvents          int     `json:"totalEvents"`
	TotalGames           int     `json:"totalGames"`
	TotalGamesInProgress int     `json:"totalGamesInProgress"`
	Dates                []*Date `json:"dates"`
}

type Date struct {
	Date                 string  `json:"date"`
	TotalItems           int     `json:"totalItems"`
	TotalEvents          int     `json:"totalEvents"`
	TotalGames           int     `json:"totalGames"`
	TotalGamesInProgress int     `json:"totalGamesInProgress"`
	Games                []*Game `json:"games"`
}

type Record struct {
	Wins   int    `json:"wins"`
	Losses int    `json:"losses"`
	Pct    string `json:"pct"`
}

type Venue struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	link string `json:"link"`
}

type League struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}

type Division struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}

type Sport struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}

func GetSchedule(date time.Time) (*Schedule, error) {
	req, err := http.NewRequest("GET", scheduleEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create schedule request: %s", err)
	}
	q := req.URL.Query()
	q.Add("date", date.Format(DateFormat))
	q.Add("sportId", sportID)
	req.URL.RawQuery = q.Encode()

	log.Printf("sending request: %s\n", req.URL.String())

	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := netClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule info: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get schedule: Status: %v, Error: %s", resp.StatusCode, resp.Body)
	}

	var schedule Schedule
	if err := json.NewDecoder(resp.Body).Decode(&schedule); err != nil {
		return nil, fmt.Errorf("error parsing schedule info: %s", err)
	}

	return &schedule, nil
}
