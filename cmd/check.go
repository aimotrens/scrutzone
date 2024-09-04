package cmd

// Represents the interface for a check
type ICheck interface {
	Validate() error
	SetDefaults(*MetaCheck)
	Run() func(*Notification)
}
