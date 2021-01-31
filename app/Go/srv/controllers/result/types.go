package result_controller

type NewResult struct {
	CheckpointID string `json:"checkpoint_id"`
	SportsmenID  string `json:"sportsmen_id"`
	EventStateID string `json:"event_state_id"`
	Time         string `json:"time"`
}
