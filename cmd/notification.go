package cmd

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/aimotrens/scrutzone/version"
)

// Represents the notification configuration
type Notification struct {
	DefaultTarget NotifyTarget       `yaml:"defaultTarget"`
	Targets       map[string]*Target `yaml:"targets"`
}

// Validates the notification configuration
func (n *Notification) Validate() error {
	if n.DefaultTarget == "" {
		return fmt.Errorf("defaultTarget is required")
	}

	if len(n.Targets) == 0 {
		return fmt.Errorf("targets required")
	}

	for k, v := range n.Targets {
		if v.Email == nil /* && <check all supported target types...> */ {
			return fmt.Errorf("no target type specified %s", k)
		}
	}

	return nil
}

// Sends a notification to the specified targets
func (n *Notification) Notify(targets []NotifyTarget, subject, text string) {
	if len(targets) == 0 {
		targets = []NotifyTarget{n.DefaultTarget}
	}

	fmt.Println("Notifying: ", targets)

	cd, ver := version.BuildInfo()
	text += "\r\n\r\n" + "----------------------------------------\r\n"
	text += fmt.Sprintf("scrutzone version: %s\r\n", ver)
	text += fmt.Sprintf("compiled at: %s\r\n", cd)

	for _, t := range targets {
		if g, ok := n.Targets[string(t)]; !ok {
			fmt.Println("Target not found: ", t)
			continue
		} else {
			if g.Email != nil {
				g.Email.sendMail(subject, text)
			} else {
				fmt.Println("Notification target type not specified")
			}
		}
	}
}

// Represents a notification target
type Target struct {
	Name  string `yaml:"name"`
	Email *Email `yaml:"email"`
}

// Represents an email configuration for notifications
type Email struct {
	Enabled bool     `yaml:"enabled"`
	Account *Account `yaml:"account"`
	Host    string   `yaml:"host"`
	Port    int      `yaml:"port"`
	From    string   `yaml:"from"`
	To      []string `yaml:"to"`
}

// Represents an email account
type Account struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// Represents a notification target name
type NotifyTarget string

// Sends a notification email
// If subject is empty, it will default to "scrutzone Notification"
func (e *Email) sendMail(subject, body string) {
	if subject == "" {
		subject = "scrutzone Notification"
	}

	sb := new(strings.Builder)
	sb.WriteString("Subject: " + subject + "\r\n")
	sb.WriteString("\r\n")
	sb.WriteString(body)

	var auth smtp.Auth

	if e.Account != nil {
		auth = smtp.PlainAuth("", e.Account.Username, e.Account.Password, e.Host)
	}

	err := smtp.SendMail(fmt.Sprintf("%s:%d", e.Host, e.Port), auth, e.From, e.To, []byte(sb.String()))
	if err != nil {
		fmt.Println(err)
		return
	}
}
