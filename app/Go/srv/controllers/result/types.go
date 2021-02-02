package result_controller

type NewResultRequest struct {
	CheckpointID string `json:"checkpoint_id"`
	SportsmenID  string `json:"sportsmen_id"`
	Time         int64  `json:"time_start"`
}

type FinishRequest struct {
	CheckpointID string `json:"checkpoint_id"`
	SportsmenID  string `json:"sportsmen_id"`
	Time         int64  `json:"time_finish"`
}
