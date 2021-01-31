package sportsmen_controller

type NewSportsmenRequest struct {
	StartNumber uint32 `json:"start_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
}
