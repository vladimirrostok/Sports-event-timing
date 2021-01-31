package dashboard_controller

type ResultMessage struct {
	ID                   string `json:"id"`
	SportsmenStartNumber string `json:"start_number"`
	SportsmenName        string `json:"name"`
	TimeStart            string `json:"time_start"`
}

type FinishMessage struct {
	ID        string `json:"id"`
	TimeStart string `json:"time_finish"`
}
