package tasmota

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/federicoleon/go-httpclient/core"
	"github.com/federicoleon/go-httpclient/gohttp"
)

// Client is a client for communication with a Tasmota-powered device.
type Client struct {
	device Device
}

// NewClient returns a new Client that is configured to be talking to the given Device.
func NewClient(device Device) *Client {
	return &Client{
		device: device,
	}
}

// SetPowerOn turns the power on for the configured device.
func (x *Client) SetPowerOn() error {
	return x.SetPower(true)
}

// SetPowerOff turns the power off for the configured device.
func (x *Client) SetPowerOff() error {
	return x.SetPower(false)
}

// SetPower sets the power state for the configured device.
func (x *Client) SetPower(on bool) error {
	val := "ON"
	if !on {
		val = "OFF"
	}
	response, err := x.get("Power " + val)
	if err != nil {
		return err
	}
	if response.String() != fmt.Sprintf(`{"POWER":"%s"}`, val) {
		return fmt.Errorf("invalid")
	}
	return nil
}

// IsPowerOn retrieves if the power is on for the configured device.
func (x *Client) IsPowerOn() (bool, error) {
	response, err := x.get("Power")
	if err != nil {
		return false, err
	}
	return response.String() == `{"POWER":"ON"}`, nil
}

// GetCurrentPower returns the current power consumption of the attached devices in full Watts.
func (x *Client) GetCurrentPower() (int, error) {
	response, err := x.get("Status 10")
	if err != nil {
		return 0, err
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(response.String()), &data)
	if err != nil {
		return 0, err
	}

	return x.getIntByPath(data, "StatusSNS.ENERGY.Power")
}

func (x *Client) httpClient() gohttp.Client {
	return gohttp.NewBuilder().
		SetConnectionTimeout(10 * time.Second).
		SetResponseTimeout(20 * time.Second).
		SetUserAgent("go-tasmota").
		Build()
}

func (x *Client) get(cmd string) (*core.Response, error) {
	apiURL := x.device.GetAPIBaseURL() + url.QueryEscape(cmd)
	return x.httpClient().Get(apiURL)
}

func (x *Client) getIntByPath(data map[string]interface{}, path string) (int, error) {
	keys := strings.Split(path, ".")

	var (
		value any
		ok    bool
	)
	for _, key := range keys {
		value, ok = data[key]
		if !ok {
			return 0, fmt.Errorf("key '%s' not found", key)
		}

		data, ok = value.(map[string]interface{})
		if !ok {
			break
		}
	}

	vInt, ok := value.(int)
	if ok {
		return vInt, nil
	}

	vFloat, ok := value.(float64)
	if ok {
		return int(vFloat), nil
	}

	return 0, fmt.Errorf("value is not an integer")
}
