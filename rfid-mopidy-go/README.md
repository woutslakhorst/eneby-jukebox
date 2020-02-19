I just the script from https://github.com/balag3/RFID_reader to find the correct device.

`lsusb` also comes in handy for the Mac

In my case for the raspberry:

```
('event: /dev/input/event0', 'name: Sycreader USB Reader', 'hardware: usb-0000:01:00.0-1.4/input0')
``` 

I develop on a Mac, so this can't be used. Using the lib from: `github.com/jpoirier/gousb/usb` we can list all devices with:

```go
    ctx := usb.NewContext()
	defer func() {
		if err := ctx.Close(); err != nil {
			fmt.Printf("failed to close usb context: %s\n", err.Error())
		}
	}()

	//ctx.Debug(1)

	// ListDevices is used to find the devices to open.
	ctx.ListDevices(
		func(desc *usb.Descriptor) bool {
			fmt.Printf("Vendor: %s, Product: %s (%v)\n", desc.Vendor, desc.Product, desc)
			return false
		})
```

The vendor we're looking for is `ffff` and the product is `0035`.
This route gives a bad access from usb C lib....

Trying with google gousb lib 
```
go get -v github.com/google/gousb/lsusb
```

code is pretty similar:

```go
    // Only one context should be needed for an application.  It should always be closed.
	ctx := gousb.NewContext()
	defer ctx.Close()

	// Iterate through available Devices, finding all that match a known VID/PID.
	vid, pid := gousb.ID(0xffff), gousb.ID(0x0035)
	devs, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		// this function is called for every device present.
		// Returning true means the device should be opened.
		return desc.Vendor == vid && desc.Product == pid
	})
	// All returned devices are now open and will need to be closed.
	for _, d := range devs {
		defer d.Close()
	}
	if err != nil {
		log.Fatalf("OpenDevices(): %v", err)
	}
	if len(devs) == 0 {
		log.Fatalf("no devices found matching VID %s and PID %s", vid, pid)
	}
```

Also doesn't work. Might have to do with the fact that my Mac has claimed the device.
Let's try a HID lib.

https://github.com/karalabe/hid seems to do what I want.
Combined with sudo it's a win for the MAC. The Read command returns keyboard codes, so some conversion was required.
Also had to give iTerm access to log keys...

For linux, I'll need to add the device to som udev list.... have to look it up

For some reason I have to write back 3 bytes to the HID device to *reset* the device, otherwise it'll return different codes for the Mac

## Cross compile

didn't work

ended up installing go on the raspberry and:

```
go get github.com/woutslakhorst/eneby-jukebox/rfid-mopidy-go
```

as a service via https://www.dexterindustries.com/howto/run-a-program-on-your-raspberry-pi-at-startup/


