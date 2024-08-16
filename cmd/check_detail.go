package cmd

// Represents the interface for a check detail
type ICheckDetail interface {
	Validate() error
	SetDefaults(*Check)
	Run() func(*Notification)
}
