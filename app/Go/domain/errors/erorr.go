package domain_errors

type (
	// InvalidVersion signifies an entity with specified version does not exist in the system.
	InvalidVersion struct{}

	// InvalidVersion signifies an entity with specified state does not exist in the system.
	StateConflict struct{}
)

func (err InvalidVersion) Error() string {
	return "Invalid version"
}

func (err StateConflict) Error() string {
	return "Invalid state"
}
