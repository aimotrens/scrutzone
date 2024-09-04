package cmd

type CheckError struct {
	checkName string
	err       error
}

func (c *CheckError) Error() string {
	return c.checkName + ": " + c.err.Error()
}
