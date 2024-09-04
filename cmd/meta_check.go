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

	previousState CheckState
	currentState  CheckState
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

	// we suggest that the check is OK by default to avoid sending unnecessary notifications on startup
	mc.previousState = CheckOk
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

// Returns the check based on the check type
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
	defer mc.RollState()

	fmt.Printf("Running check %s\n", mc.Name)

	check, err := mc.getCheck()
	if err != nil {
		mc.SetState(CheckFailed)
		return mc.DefaultNotifyFailedFunc(OnStateChanged, mc.NewError(err))()
	}

	s, n := check.Run()
	mc.SetState(s)

	return n()
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

// Sets the previous state of the check
func (mc *MetaCheck) SetState(s CheckState) {
	mc.currentState = s
}

// Returns a default notification function for errors
// If cond() returns true, it returns a func that sends a notification
// If cond() returns false and the current check state is CheckFailed, it returns a func that prints only a log message
// Otherwise, it returns nil
func (mc *MetaCheck) DefaultNotifyFailedFunc(cond func(mc *MetaCheck) bool, err *CheckError) NotifyFuncSwitch {
	return func() func(*Notification) {
		n := func(n *Notification) {
			n.Notify(
				mc.NotifyTargets,
				fmt.Sprintf("scrutzone | check failure for %s", mc.Name),
				err.Error(),
			)
		}
		l := func(n *Notification) {
			fmt.Println("Notification already sent for check ", mc.Name)
		}

		if cond(mc) {
			return n
		}

		// only log if the notification has already been sent
		if mc.previousState == CheckFailed {
			return l
		}

		return nil
	}
}

// Returns a default notification function for OK messages
// If cond() returns true, it returns a func that sends a notification
// Otherwise, it returns nil
func (mc *MetaCheck) DefaultNotifyOkFunc(cond func(mc *MetaCheck) bool) NotifyFuncSwitch {
	return func() func(*Notification) {
		n := func(n *Notification) {
			n.Notify(mc.NotifyTargets, "scrutzone Check OK", fmt.Sprintf("Check %s OK", mc.Name))
		}

		if cond(mc) {
			return n
		}

		return nil
	}
}

// Checks if the state has changed
func (mc *MetaCheck) HasStateChanged() bool {
	return mc.previousState != mc.currentState
}

// Rolls the state of the check
func (m *MetaCheck) RollState() {
	m.previousState = m.currentState
}

// Creates a new error for the check
func (mc *MetaCheck) NewError(err error) *CheckError {
	return &CheckError{
		checkName: mc.Name,
		err:       err,
	}
}

// Activates the notification if the state has changed
func OnStateChanged(mc *MetaCheck) bool {
	return mc.HasStateChanged()
}

// Always activates the notification
func Always(mc *MetaCheck) bool {
	return true
}

// Never activates the notification
func Never(mc *MetaCheck) bool {
	return false
}
