package tcpcheck

import (
	"fmt"
	"net"

	"github.com/aimotrens/scrutzone/cmd"
)

func init() {
	cmd.RegisterCheckDetail("tcp", NewTcpCheck)
}

// Represents the configuration for a TCP check
type TcpCheck struct {
	check *cmd.Check

	Port int `yaml:"port"`
}

// Creates a new TCP check
func NewTcpCheck(c *cmd.Check) cmd.ICheckDetail {
	return &TcpCheck{
		check: c,
	}
}

// Sets the default values for the TCP check
func (t *TcpCheck) SetDefaults(c *cmd.Check) {
}

// Validates the TCP check configuration
func (t *TcpCheck) Validate() error {
	return nil
}

// Runs the TCP check
func (t *TcpCheck) Run() func(*cmd.Notification) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", t.check.Address, t.Port))
	if err != nil {
		return t.check.DefaultNotifyFunc(err)
	}
	defer conn.Close()

	return nil
}
