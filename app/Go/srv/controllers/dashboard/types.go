package dashboard_controller

type ResultMessage struct {
	ID                   string `json:"id"`
	SportsmenStartNumber uint32 `json:"start_number"`
	SportsmenName        string `json:"name"`
	TimeStart            int64  `json:"time_start"`
	TimeFinish           *int64 `json:"time_finish"`
}

type UnfinishedResultMessage struct {
	ID                   string `json:"id"`
	SportsmenStartNumber uint32 `json:"start_number"`
	SportsmenName        string `json:"name"`
	TimeStart            int64  `json:"time_start"`
}

type FinishedResultMessage struct {
	ID                   string `json:"id"`
	SportsmenStartNumber uint32 `json:"start_number"`
	SportsmenName        string `json:"name"`
	TimeFinish           int64  `json:"time_finish"`
}
