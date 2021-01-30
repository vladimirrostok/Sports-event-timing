package main

type Result struct {
	ID           string `json:"id"`
	CheckpointID string `json:"checkpoint_id"`
	SportsmenID  string `json:"sportsmen_id"`
	EventStateID string `json:"event_state_id"`
	Time         string `json:"time"`
}

func NewResult(id, checkpointID, sportsmenID, eventStateID, time string) *Result {
	return &Result{
		ID:           id,
		CheckpointID: checkpointID,
		SportsmenID:  sportsmenID,
		EventStateID: eventStateID,
		Time:         time,
	}
}
