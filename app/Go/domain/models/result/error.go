package result

type (
	// NotFound signifies a result is not found.
	NotFound struct{}
)

func (err NotFound) Error() string {
	return "Result does not exist"
}
