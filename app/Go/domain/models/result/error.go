package result

type (
	// AlreadyExists signifies an result with already exists in the system.
	AlreadyExists struct{}

	// AlreadyFinished signifies an result with finish time set already.
	AlreadyFinished struct{}

	// NotFound signifies a result is not found.
	NotFound struct{}
)

func (err AlreadyExists) Error() string {
	return "Result already exists"
}

func (err NotFound) Error() string {
	return "Result does not exist"
}

func (err AlreadyFinished) Error() string {
	return "Result has finish time already"
}
