package sportsmen

type (
	// NotFound signifies a sportsmen is not found.
	NotFound struct{}
)

func (err NotFound) Error() string {
	return "Sportsmen does not exist"
}
