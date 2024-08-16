package cmd

// Represents the global main configuration
type Config struct {
	StartupNotification *StartupNotification `yaml:"startupNotification"`
	CheckConfigDir      string               `yaml:"checkConfigDir"`
	Notification        *Notification        `yaml:"notification"`
	CheckDefaults       *CheckDefaults       `yaml:"checkDefaults"`

	Checks []*Check
}
