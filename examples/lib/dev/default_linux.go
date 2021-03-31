package dev

import (
	"github.com/wmh11112345/ble"
	"github.com/wmh11112345/ble/linux"
)

// DefaultDevice ...
func DefaultDevice(opts ...ble.Option) (d ble.Device, err error) {
	return linux.NewDevice(opts...)
}
