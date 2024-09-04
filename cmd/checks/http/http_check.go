package httpcheck

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/aimotrens/scrutzone/cmd"
)

func init() {
	cmd.RegisterCheckDetail("http", NewHttpCheck)
}

// Represents the configuration for an HTTP check
type HttpCheck struct {
	check *cmd.MetaCheck

	Hostname     string     `yaml:"hostname"`
	Port         int        `yaml:"port"`
	Scheme       HttpScheme `yaml:"scheme`
	Path         string     `yaml:"path"`
	Username     string     `yaml:"username"`
	Password     string     `yaml:"password"`
	ExpectedCode uint       `yaml:"expectedCode"`
}

type HttpScheme string

func (hs HttpScheme) Validate() error {
	allowedSchemes := []string{"http", "https"}
	if !slices.Contains(allowedSchemes, string(hs)) {
		return fmt.Errorf("unsupported scheme: %s", hs)
	}
	return nil
}

// Creates a new HTTP check
func NewHttpCheck(c *cmd.MetaCheck) cmd.ICheck {
	return &HttpCheck{
		check: c,
	}
}

// Sets the default values for the HTTP check
func (h *HttpCheck) SetDefaults(c *cmd.MetaCheck) {
	if h.Hostname == "" {
		h.Hostname = c.Address
	}

	if h.Scheme == "" {
		h.Scheme = HttpScheme("http")
	}

	if h.Port == 0 {
		switch h.Scheme {
		case "http":
			h.Port = 80
		case "https":
			h.Port = 443
		}
	}

	if h.Path == "" {
		h.Path = "/"
	}

	if h.ExpectedCode == 0 {
		h.ExpectedCode = 200
	}
}

// Validates the HTTP check configuration
func (h *HttpCheck) Validate() error {
	if h.ExpectedCode < 100 || h.ExpectedCode > 599 {
		return fmt.Errorf("expected code must be between 100 and 599")
	}

	if err := h.Scheme.Validate(); err != nil {
		return err
	}

	return nil
}

// Runs the HTTP check
func (h *HttpCheck) Run() (cmd.CheckState, cmd.NotifyFuncSwitch) {
	url := fmt.Sprintf("%s://%s:%d%s", h.Scheme, h.Hostname, h.Port, h.Path)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return cmd.CheckFailed, h.check.DefaultNotifyFailedFunc(cmd.OnStateChanged, h.check.NewError(err))
	}

	if h.Username != "" || h.Password != "" {
		req.SetBasicAuth(h.Username, h.Password)
	}

	httpClient := http.DefaultClient
	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return cmd.CheckFailed, h.check.DefaultNotifyFailedFunc(cmd.OnStateChanged, h.check.NewError(err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != int(h.ExpectedCode) {
		err := fmt.Errorf("expected code %d, got %d", h.ExpectedCode, resp.StatusCode)
		return cmd.CheckFailed, h.check.DefaultNotifyFailedFunc(cmd.OnStateChanged, h.check.NewError(err))
	}

	return cmd.CheckOk, h.check.DefaultNotifyOkFunc(cmd.OnStateChanged)
}
