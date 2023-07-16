package tasmota

import (
	"fmt"
	"strings"
)

// Device is a Tasmota-powered device
type Device struct {
	Hostname string
}

// GetAPIBaseURL returns the base URL to be used to send commands to this Tasmota device.
func (x Device) GetAPIBaseURL() string {
	hostname := strings.TrimPrefix(x.Hostname, "http://")
	return fmt.Sprintf("http://%s/cm?cmnd=", hostname)
}
