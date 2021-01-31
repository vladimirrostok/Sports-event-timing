package dashboard_controller

type ResultMessage struct {
	ID                   string `json:"id"`
	SportsmenStartNumber string `json:"sportsmen_id"`
	SportsmenName        string `json:"sportsmen_name"`
	TimeStart            string `json:"time_start"`
}

type FinishMessage struct {
	ID        string `json:"id"`
	TimeStart string `json:"time_finish"`
}
