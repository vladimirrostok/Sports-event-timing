package checkpoint_controller

type NewCheckpointRequest struct {
	Name string `json:"name"`
}

type CreatedResponse struct {
	ID string `json:"id"`
}
