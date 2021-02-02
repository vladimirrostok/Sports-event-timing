package checkpoint

type (
	// NotFound signifies a checkpoint is not found.
	NotFound struct{}
)

func (err NotFound) Error() string {
	return "Checkpoint does not exist"
}
