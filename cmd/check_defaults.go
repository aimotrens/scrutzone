package cmd

import "errors"

// CheckDefaults represents the global default values for checks
type CheckDefaults struct {
	NotifyTargets []NotifyTarget `yaml:"notifyTargets"`
	Interval      *int           `yaml:"interval"`
	Timeout       *int           `yaml:"timeout"`
}

func (c *CheckDefaults) Validate() error {
	errs := []error{}

	if c.Interval == nil {
		c.Interval = new(int)
		*c.Interval = 60
	} else if *c.Interval <= 0 {
		errs = append(errs, errors.New("interval must be greater than 0"))
	}

	if c.Timeout == nil {
		c.Timeout = new(int)
		*c.Timeout = 5
	} else if *c.Timeout < 0 {
		errs = append(errs, errors.New("timeout must be greater than or equal to 0"))
	}

	return errors.Join(errs...)
}
