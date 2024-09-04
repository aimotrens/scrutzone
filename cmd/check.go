package cmd

type NotifyFuncSwitch func() func(*Notification)

// Represents the interface for a check
type ICheck interface {
	Validate() error
	SetDefaults(*MetaCheck)
	Run() (CheckState, NotifyFuncSwitch)
}
