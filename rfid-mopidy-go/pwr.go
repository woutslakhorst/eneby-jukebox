package main

import (
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/woutslakhorst/go-rpio"
)

// monitor BCM03 (physical pin 5)
func monitorPin5() error {
	// open memory space
	if err := rpio.Open(); err != nil {
		return err
	}

	pin := rpio.Pin(3)
	pin.Input()
	pin.PullUp()
	pin.Detect(rpio.FallEdge)

	log.Print("Starting monitoring of BCM03")

	go func() {
		defer rpio.Close()
		var detected bool

		for !detected {
			detected = pin.EdgeDetected()
			time.Sleep(100 * time.Millisecond)
		}

		pin.Detect(rpio.NoEdge)

		// halt
		cmd := exec.Command("halt")
		log.Print("Stopping")
		if err := cmd.Run(); err != nil {
			log.Printf("failed to issue halt command: %s", err.Error())
		}
		os.Exit(0)
	}()

	return nil
}
