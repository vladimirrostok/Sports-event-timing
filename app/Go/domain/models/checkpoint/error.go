package checkpoint

type (
	// AlreadyExists signifies a checkpoint with specified name already exists in the system.
	AlreadyExists struct{}

	// NotFound signifies a checkpoint is not found.
	NotFound struct{}
)

func (err AlreadyExists) Error() string {
	return "Checkpoint already exists"
}

func (err NotFound) Error() string {
	return "Checkpoint does not exist"
}
