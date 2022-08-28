package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/kevinjpickard/mlb-personalized-schedule/schedule"
)

const (
	serverPort = ":8080"
)

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	switch {
	case r.Method == http.MethodGet:
		queryParams, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("malformed request: %s", err)))
		}

		date := queryParams.Get("date")
		if date == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("date is required"))
			return
		}
		parsedDate, err := time.Parse(schedule.DateFormat, date)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("date not recognized, expected format YYYY-MM-DD, got: %s", date)))
			return
		}

		teamID := queryParams.Get("teamId")
		if teamID == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("teamId is required"))
			return
		}
		teamIDNumber, err := strconv.Atoi(teamID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("teamId not recognized: %s", err)))
		}

		schedule, err := schedule.GetSchedule(parsedDate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("could not get schedule: %s", err)))
			return
		}

		schedule.SortScheduleByTeam(teamIDNumber)
		resp, err := json.Marshal(schedule)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("could not parse schedule: %s", err)))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
		return

	default:
		http.NotFound(w, r)
		return
	}
}

func main() {
	http.HandleFunc("/", ServeHTTP)
	log.Fatal(http.ListenAndServe(serverPort, nil))
}
