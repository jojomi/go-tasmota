package tasmota

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPower(t *testing.T) {
	a := assert.New(t)
	c := NewClient(Device{
		Hostname: "http://midnight-blue",
	})

	err := c.SetPowerOn()
	a.Nil(err)
	on, err := c.IsPowerOn()
	a.Nil(err)
	a.True(on)

	time.Sleep(5 * time.Second)

	pow, err := c.GetCurrentPower()
	a.Nil(err)
	a.True(pow > 0.0)

	err = c.SetPowerOff()
	a.Nil(err)
	on, err = c.IsPowerOn()
	a.Nil(err)
	a.False(on)

	time.Sleep(2 * time.Second)

	pow, err = c.GetCurrentPower()
	a.Nil(err)
	a.Equal(0.0, pow)
}
