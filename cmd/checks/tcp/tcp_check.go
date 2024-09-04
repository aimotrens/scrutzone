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
	check *cmd.MetaCheck

	Port int `yaml:"port"`
}

// Creates a new TCP check
func NewTcpCheck(c *cmd.MetaCheck) cmd.ICheck {
	return &TcpCheck{
		check: c,
	}
}

// Sets the default values for the TCP check
func (t *TcpCheck) SetDefaults(c *cmd.MetaCheck) {
}

// Validates the TCP check configuration
func (t *TcpCheck) Validate() error {
	return nil
}

// Runs the TCP check
func (t *TcpCheck) Run() (cmd.CheckState, cmd.NotifyFuncSwitch) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", t.check.Address, t.Port))
	if err != nil {
		return cmd.CheckFailed, t.check.DefaultNotifyFailedFunc(cmd.OnStateChanged, t.check.NewError(err))
	}
	defer conn.Close()

	return cmd.CheckOk, t.check.DefaultNotifyOkFunc(cmd.OnStateChanged)
}
