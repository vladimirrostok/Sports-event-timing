package eventstate

type (
	// AlreadyExists signifies a state with specified name already exists in the system.
	AlreadyExists struct{}

	// NotFound signifies a state is not found.
	NotFound struct{}
)

func (err AlreadyExists) Error() string {
	return "Event state already exists"
}

func (err NotFound) Error() string {
	return "Event state does not exist"
}
