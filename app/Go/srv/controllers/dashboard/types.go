package dashboard_controller

type ResultMessage struct {
	ID                   string `json:"id"`
	SportsmenStartNumber string `json:"start_number"`
	SportsmenName        string `json:"name"`
	TimeStart            string `json:"time_start"`
	TimeFinish           string `json:"time_finish"`
}

type UnfinishedResultMessage struct {
	ID                   string `json:"id"`
	SportsmenStartNumber string `json:"start_number"`
	SportsmenName        string `json:"name"`
	TimeStart            string `json:"time_start"`
}

type FinishedResultMessage struct {
	ID                   string `json:"id"`
	SportsmenStartNumber string `json:"start_number"`
	SportsmenName        string `json:"name"`
	TimeFinish           string `json:"time_finish"`
}
