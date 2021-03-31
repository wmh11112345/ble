package dev

import (
	"github.com/wmh11112345/ble"
	"github.com/wmh11112345/ble/darwin"
)

// DefaultDevice ...
func DefaultDevice(opts ...ble.Option) (d ble.Device, err error) {
	return darwin.NewDevice(opts...)
}
