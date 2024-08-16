package cmd

type CheckError struct {
	checkName string
	err       error
}

func NewCheckError(name string, err error) error {
	return &CheckError{
		checkName: name,
		err:       err,
	}
}

func (c *CheckError) Error() string {
	return c.checkName + ": " + c.err.Error()
}
