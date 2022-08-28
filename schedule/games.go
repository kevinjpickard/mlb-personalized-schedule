package schedule

import "log"

const (
	GameStatus_Scheduled  = "S"
	GameStatus_Final      = "F"
	GameStatus_InProgress = "I"
	GameStatus_Postponed  = "D" // Darn COVID
)

// dates are strings for now, if there is enough time I will write a marshaller/unmarshaller to directly translate
// these to dateTime objects for easier handling.
type Game struct {
	GamePk                 int         `json:"gamePk"`
	Link                   string      `json:"link"`
	GameType               string      `json:"gameType"`
	Season                 string      `json:"season"`
	GameDate               string      `json:"gameDate"`
	OfficialDate           string      `json:"officialDate"`
	RescheduleDate         string      `json:"rescheduleDate"`     // Only if postponed
	RescheduleGameDate     string      `json:"rescheduleGameDate"` // Only if postponed
	Status                 *GameStatus `json:"status"`
	Teams                  *Side       `json:"teams"`
	Venue                  *Venue      `json:"venue"`
	Content                *Content    `json:"content"`
	IsTie                  bool        `json:"isTie"`
	GameNumber             int         `json:"gameNumber"`
	PublicFacing           bool        `json:"publicFacing"`
	DoubleHeader           string      `json:"doubleHeader"`
	GamedayType            string      `json:"gamedayType"`
	TieBreaker             string      `json:"tieBreaker"`
	CalendarEventID        string      `json:"calendarEventID"`
	SeasonDisplay          string      `json:"seasonDisplay"`
	dayNight               string      `json:"dayNight"`
	ScheduledInnings       uint        `json:"scheduledInnings"`
	ReverseHomeAwayStatus  bool        `json:"reverseHomeAwayStatus"`
	InningBreakLength      int         `json:"inningBreakLength"`
	GamesInSeries          int         `json:"gamesInSeries"`
	SeriesGameNumber       int         `json:"seriesGameNumber"`
	SeriesDescription      string      `json:"seriesDescription"`
	RecordSource           string      `json:"recordSource"`
	IfNecessary            string      `json:"ifNecessary"`
	IfNecessaryDescription string      `json:"ifNecessaryDescription"`
}

type Content struct {
	Link string `json:"link"`
}

type TeamInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}

type Team struct {
	LeagueRecord *Record   `json:"leagueRecord"`
	Score        int       `json:"score"`
	Team         *TeamInfo `json:"team"`
	IsWinner     bool      `json:"isWinner"`
	SplitSquad   bool      `json:"splitSquad"`
	SeriesNumber int       `json:"seriesNumber"`
}

type Side struct {
	Away *Team `json:"away"`
	Home *Team `json:"home"`
}

type GameStatus struct {
	AbstractGameState string `json:"abstractGameState"`
	CodedGameState    string `json:"codedGameState"`
	DetailedState     string `json:"detailedState"`
	StatusCode        string `json:"statusCode"`
	StartTimeTBD      bool   `json:"startTimeTBD"`
}

// Assuming MLB API always returns a stable, chronological order (which is true in my limited experience)
func (s *Schedule) SortScheduleByTeam(teamID int) {
	log.Printf("sorting for team with ID: %d\n", teamID)

	var (
		priorityGames []*Game
		otherGames    []*Game
	)

	// The endpoint does not support more than a single day, so the s.Dates slice will never have more than 1 element.
	for _, game := range s.Dates[0].Games {
		if game.Teams.Home.Team.ID == teamID || game.Teams.Away.Team.ID == teamID {
			priorityGames = append(priorityGames, game)
			// Triple-Headers are prohibited under the CBA, so this slice will never contain more than 2 games
			if len(priorityGames) > 1 {
				log.Println("Detected DoubleHeader")
				// Only need to reorder if the second game is in progress
				if priorityGames[1].Status.CodedGameState == GameStatus_InProgress {
					log.Println("DoubleHeader in progress, reordering")
					priorityGames = []*Game{priorityGames[1], priorityGames[0]}
				}
			}
		} else {
			otherGames = append(otherGames, game)
		}
	}

	s.Dates[0].Games = append(priorityGames, otherGames...)
}
