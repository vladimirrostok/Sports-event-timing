package main

type Result struct {
	ID                   string `json:"id"`
	SportsmenStartNumber string `json:"sportsmen_id"`
	SportsmenName        string `json:"sportsmen_name"`
	Time                 string `json:"finish_time"`
}

func NewResult(id, sportsmenStartNumber, sportsmenName, time string) *Result {
	return &Result{
		ID:                   id,
		SportsmenStartNumber: sportsmenStartNumber,
		SportsmenName:        sportsmenName,
		Time:                 time,
	}
}
