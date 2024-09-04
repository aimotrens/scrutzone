package cmd

// Represents the state of a check
type CheckState bool

const (
	// CheckOk represents a successful check
	CheckOk CheckState = true

	// CheckFailed represents a failed check
	CheckFailed CheckState = false
)
