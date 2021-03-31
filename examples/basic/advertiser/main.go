package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/wmh11112345/ble/linux"

	"github.com/wmh11112345/ble"
	"github.com/wmh11112345/ble/examples/lib/dev"
	"github.com/wmh11112345/ble/linux/hci/cmd"
	"github.com/pkg/errors"
)

var (
	device = flag.String("device", "default", "implementation of ble")
	du     = flag.Duration("du", 10*time.Second, "advertising duration, 0 for indefinitely")
)

func updateLinuxAdvParam(d *linux.Device) error {
	if err := d.HCI.Send(&cmd.LESetAdvertisingParameters{
		AdvertisingIntervalMin:  80,        // 0x0020 - 0x4000; N * 0.625 msec
		AdvertisingIntervalMax:  80,        //80,        // 0x0020 - 0x4000; N * 0.625 msec
		AdvertisingType:         0x00,      // 00: ADV_IND, 0x01: DIRECT(HIGH), 0x02: SCAN, 0x03: NONCONN, 0x04: DIRECT(LOW)
		OwnAddressType:          0x00,      // 0x00: public, 0x01: random
		DirectAddressType:       0x00,      // 0x00: public, 0x01: random
		DirectAddress:           [6]byte{}, // Public or Random Address of the Device to be connected
		AdvertisingChannelMap:   0x01,      // 0x07 0x01: ch37, 0x2: ch38, 0x4: ch39
		AdvertisingFilterPolicy: 0x00,
	}, nil); err != nil {
		return errors.Wrap(err, "can't set advertising param")
	}
	LEReadAdvertisingChannelTxPowerRP := cmd.LEReadAdvertisingChannelTxPowerRP{}
	d.HCI.Send(&cmd.LEReadAdvertisingChannelTxPower{}, &LEReadAdvertisingChannelTxPowerRP)
	fmt.Printf("TransmitPowerLevel:%d", int(LEReadAdvertisingChannelTxPowerRP.TransmitPowerLevel))
	return nil
}

func main() {
	flag.Parse()

	d, err := dev.NewDevice("default")
	if err != nil {
		log.Fatalf("can't new device : %s", err)
	}
	ble.SetDefaultDevice(d)
	// Advertise for specified durantion, or until interrupted by user.
	fmt.Printf("Advertising for %s...\n", *du)

	//chkErr(ble.AdvertiseNameAndServices(ctx, "wmh"))
	//var pwr int8 = -20
	if dev, ok := d.(*linux.Device); ok {
		if err := errors.Wrap(updateLinuxAdvParam(dev), "can't update hci parameters"); err != nil {
			panic(err)
		}
	}
	for {
		ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), *du))
		chkErr(ble.AdvertiseMfgDataDIY(ctx, 0x55cc, []byte{0x1, 0x2, 0x3, 0x4}))
		time.Sleep(10 * time.Second)
	}

	//chkErr(ble.AdvertiseIBeacon(ctx, []byte{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x4B, 0xEE, 0x95, 0xF7, 0xD8, 0xCC, 0x64, 0xA8, 0x63, 0xB5}, 1, 2, pwr))
}

func chkErr(err error) {
	switch errors.Cause(err) {
	case nil:
	case context.DeadlineExceeded:
		fmt.Printf("done\n")
	case context.Canceled:
		fmt.Printf("canceled\n")
	default:
		log.Fatalf(err.Error())
	}
}
