package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

var checkDetails = make(map[string]func(c *Check) ICheckDetail)

func RegisterCheckDetail(name string, checkConstructor func(c *Check) ICheckDetail) {
	checkDetails[name] = checkConstructor
}

// Represents a check configuration
type Check struct {
	Name          string
	Address       string         `yaml:"address"`
	Interval      int            `yaml:"interval"`
	Timeout       int            `yaml:"timeout"`
	Type          string         `yaml:"type"`
	Config        map[string]any `yaml:"config"`
	NotifyTargets []NotifyTarget `yaml:"notifyTargets"`

	notificationSent bool
}

// Sets the default values for the check
func (c *Check) SetDefaults(def *CheckDefaults) {
	if c.NotifyTargets == nil {
		c.NotifyTargets = def.NotifyTargets
	}

	if c.Timeout == 0 {
		c.Timeout = *def.Timeout
	}

	if c.Interval == 0 {
		c.Interval = *def.Interval
	}
}

// Validates the check configuration
func (c *Check) Validate() error {
	errs := []error{}

	if c.Interval <= 0 {
		errs = append(errs, errors.New("interval is required and must be greater than 0"))
	}

	if c.Timeout < 0 {
		errs = append(errs, errors.New("timeout is required and must be greater than or equal to 0"))
	}

	if c.Type == "" {
		errs = append(errs, errors.New("type is required"))
	}

	if check, err := c.getCheckDetail(); err != nil {
		errs = append(errs, err)
	} else {
		errs = append(errs, check.Validate())
	}

	return errors.Join(errs...)
}

// Returns the check detail based on the check type
func (c *Check) getCheckDetail() (ICheckDetail, error) {
	data, _ := yaml.Marshal(c.Config)
	d := yaml.NewDecoder(bytes.NewReader(data))
	d.KnownFields(true)

	var err error
	var check ICheckDetail

	if constructor, ok := checkDetails[c.Type]; ok {
		check = constructor(c)
	} else {
		return nil, fmt.Errorf("unknown check type %s", c.Type)
	}

	err = d.Decode(check)

	if err == nil {
		check.SetDefaults(c)
	}

	return check, err
}

// Runs the check
func (c *Check) Run() func(*Notification) {
	fmt.Printf("Running check %s\n", c.Name)

	check, err := c.getCheckDetail()
	if err != nil {
		return c.DefaultNotifyFunc(err)
	}

	n := check.Run()
	if n == nil {
		c.ResetNotificationSent()
	}

	return n
}

// Initializes the ticker for the check
func (c *Check) InitTicker(runQueue chan<- func() func(*Notification)) {
	t := time.NewTicker(time.Duration(c.Interval) * time.Second)

	go func() {
		runQueue <- c.Run

		for range t.C {
			runQueue <- c.Run
		}
	}()
}

// Sets the notification sent flag
func (c *Check) SetNotificationSent() {
	c.notificationSent = true
}

// Resets the notification sent flag
func (c *Check) ResetNotificationSent() {
	c.notificationSent = false
}

// Returns the notification sent flag
func (c *Check) IsNotificationSent() bool {
	return c.notificationSent
}

// Returns a default notification function
// If the notification has already been sent, it returns a function that only prints to Stdout and does not send a notification
func (c *Check) DefaultNotifyFunc(err error) func(*Notification) {
	n := func(n *Notification) {
		n.Notify(c.NotifyTargets, "scrutzone Check Failure", NewCheckError(c.Name, err).Error())
	}
	l := func(n *Notification) {
		fmt.Println("Notification already sent for check ", c.Name)
	}

	if !c.IsNotificationSent() {
		c.SetNotificationSent()
		return n
	}

	return l
}
