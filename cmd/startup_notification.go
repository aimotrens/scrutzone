package cmd

// Represents the startup notification configuration
type StartupNotification struct {
	Targets []NotifyTarget `yaml:"targets"`
}
