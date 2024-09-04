package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

var checks = make(map[string]func(c *MetaCheck) ICheck)

func RegisterCheckDetail(name string, checkConstructor func(c *MetaCheck) ICheck) {
	checks[name] = checkConstructor
}

// Represents a check configuration
type MetaCheck struct {
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
func (mc *MetaCheck) SetDefaults(def *CheckDefaults) {
	if mc.NotifyTargets == nil {
		mc.NotifyTargets = def.NotifyTargets
	}

	if mc.Timeout == 0 {
		mc.Timeout = *def.Timeout
	}

	if mc.Interval == 0 {
		mc.Interval = *def.Interval
	}
}

// Validates the check configuration
func (mc *MetaCheck) Validate() error {
	errs := []error{}

	if mc.Interval <= 0 {
		errs = append(errs, errors.New("interval is required and must be greater than 0"))
	}

	if mc.Timeout < 0 {
		errs = append(errs, errors.New("timeout is required and must be greater than or equal to 0"))
	}

	if mc.Type == "" {
		errs = append(errs, errors.New("type is required"))
	}

	if check, err := mc.getCheck(); err != nil {
		errs = append(errs, err)
	} else {
		errs = append(errs, check.Validate())
	}

	return errors.Join(errs...)
}

// Returns the check detail based on the check type
func (mc *MetaCheck) getCheck() (ICheck, error) {
	data, _ := yaml.Marshal(mc.Config)
	d := yaml.NewDecoder(bytes.NewReader(data))
	d.KnownFields(true)

	var err error
	var check ICheck

	if constructor, ok := checks[mc.Type]; ok {
		check = constructor(mc)
	} else {
		return nil, fmt.Errorf("unknown check type %s", mc.Type)
	}

	err = d.Decode(check)

	if err == nil {
		check.SetDefaults(mc)
	}

	return check, err
}

// Runs the check
func (mc *MetaCheck) Run() func(*Notification) {
	fmt.Printf("Running check %s\n", mc.Name)

	check, err := mc.getCheck()
	if err != nil {
		return mc.DefaultNotifyFunc(err)
	}

	n := check.Run()
	if n == nil {
		mc.ResetNotificationSent()
	}

	return n
}

// Initializes the ticker for the check
func (mc *MetaCheck) InitTicker(runQueue chan<- func() func(*Notification)) {
	t := time.NewTicker(time.Duration(mc.Interval) * time.Second)

	go func() {
		runQueue <- mc.Run

		for range t.C {
			runQueue <- mc.Run
		}
	}()
}

// Sets the notification sent flag
func (mc *MetaCheck) SetNotificationSent() {
	mc.notificationSent = true
}

// Resets the notification sent flag
func (mc *MetaCheck) ResetNotificationSent() {
	mc.notificationSent = false
}

// Returns the notification sent flag
func (mc *MetaCheck) IsNotificationSent() bool {
	return mc.notificationSent
}

// Returns a default notification function
// If the notification has already been sent, it returns a function that only prints to Stdout and does not send a notification
func (mc *MetaCheck) DefaultNotifyFunc(err error) func(*Notification) {
	n := func(n *Notification) {
		n.Notify(mc.NotifyTargets, "scrutzone Check Failure", NewCheckError(mc.Name, err).Error())
	}
	l := func(n *Notification) {
		fmt.Println("Notification already sent for check ", mc.Name)
	}

	if !mc.IsNotificationSent() {
		mc.SetNotificationSent()
		return n
	}

	return l
}
