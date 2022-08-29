# MLB Take-Home Assignment

## Description
This is a basic implementation of a web server for sorting a Major League Baseball schedule based on a user's favorite team. It takes as input the date (in the form `2006-01-02`), the team ID (e.g. `104`) and returns the schedule of games for that day, with the specified team's games to the top of the list. If the specified team has a doubleheader, both games will be pushed to the top. If either of the games is in progress, that game will be at the top of the list, otherwise, they will be in chronological order.

This was implemented using only the Golang standard library. Of course, there are many ways to implement this, but I felt that using the Golang standard library would demonstrate how I would and implement code in general. If I had not restricted myself to using the standard library, I would have utilized Ginkgo and Gomega for testing and would have used a tool for specifying the API interfaces (such as swagger).

## Setup/Running
Install and configure [Golang](https://go.dev/doc/install). Then run the server from the root of the repo:
```
go run main.go
```
and then send a request:
```
curl -s 'localhost:8080/?date=2021-04-01&teamId=110'
```
change `date` and `teamId` accordingly. To find the ID for a specific team, make the following request:
```
curl -s https://statsapi.mlb.com/api/v1/teams\?season\=2021\&sportId\=1
```
It can be helpful to use [jq](https://stedolan.github.io/jq/download/) to filter the results. Replace the team name in the `jq` query accordingly:
```
curl -s https://statsapi.mlb.com/api/v1/teams\?season\=2021\&sportId\=1 | jq '.teams[] | select(.teamName == "Rockies")'
```

## Testing
Testing is implemented using the Golang testing package. To run:
```
go test ./...
```

## Future Developments/Improvements
As this was a take-home assignment, there are a number of changes I would make given more time. Some of these are outlined below.
* More complete and thorough testing package
* Add a benchmarking suite
* Add a docker/docker-compose layer for portability and ease of use
* Add code coverage
* Add CI/CD (typically TravisCI or CircleCI)
