package schedule_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/kevinjpickard/mlb-personalized-schedule/schedule"
)

func TestSortScheduleByTeam(t *testing.T) {
	var tests = []struct {
		teamID     int
		dateString string
	}{
		{110, "2021-04-01"},
		{160, "2021-07-13"},
		{110, "2021-09-11"},
		{117, "2021-11-02"},
		{136, "2022-03-31"},
		{137, "2022-07-21"},
	}

	for _, tt := range tests {
		path := filepath.Join("testdata", tt.dateString+".input.json")
		_, filename := filepath.Split(path)
		testname := filename[:len(filename)-len(filepath.Ext(path))]

		t.Run(testname, func(t *testing.T) {
			source, err := os.ReadFile(path)
			if err != nil {
				t.Fatal("failed to read source file:", err)
			}

			var testSchedule schedule.Schedule
			err = json.Unmarshal(source, &testSchedule)
			if err != nil {
				t.Fatal("failed to unmarshal test file:", err)
			}
			testSchedule.SortScheduleByTeam(tt.teamID)

			if testSchedule.Dates[0].Games[0].Teams.Away.Team.ID != tt.teamID &&
				testSchedule.Dates[0].Games[0].Teams.Home.Team.ID != tt.teamID {
				t.Errorf("%s not sorted correctly for team %d:", tt.dateString, tt.teamID)
			}
			if testSchedule.Dates[0].Games[0].DoubleHeader != "N" {
				if testSchedule.Dates[0].Games[1].Teams.Away.Team.ID != tt.teamID &&
					testSchedule.Dates[0].Games[1].Teams.Home.Team.ID != tt.teamID {
					t.Errorf("%s not sorted correctly for team %d:", tt.dateString, tt.teamID)
				}
			}
		})
	}
}
