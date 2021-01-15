package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/xarantolus/backtap/input"
	"github.com/xarantolus/backtap/listener"
)

func main() {
	fmt.Println("Start")

	// we use this context to make sure all processes we start etc.
	// are closed/killed correctly in case something goes wrong
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up the stream we're listening to
	fingerEvents, err := listener.LogCat(ctx)
	if err != nil {
		panic("setting up logcat listener: " + err.Error())
	}

	// Open all device files we want to write to

	// Vibrator: When we write a number as string, the devices vibrates for that amount in milliseconds
	vibratorDevice, err := os.OpenFile("/sys/devices/virtual/timed_output/vibrator/enable", os.O_WRONLY, os.ModeDevice)
	if err != nil {
		panic("cannot open vibrator device: " + err.Error())
	}
	defer vibratorDevice.Close()

	// Key: This input allows us to simulate a power button press
	keyDevice, err := os.OpenFile("/dev/input/event0", os.O_WRONLY, os.ModeDevice)
	if err != nil {
		panic("opening touch device input file: " + err.Error())
	}
	defer keyDevice.Close()

	// Touch: This input allows us to tap on the screen
	// This also works when the lockscreen is shown, so we should
	// be rather careful not to call emergency services accidentally
	touchDevice, err := os.OpenFile("/dev/input/event1", os.O_WRONLY, os.ModeDevice)
	if err != nil {
		panic("opening touch device input file: " + err.Error())
	}
	defer touchDevice.Close()

	var (
		lastTapTime   time.Time
		lastPowerTime time.Time
	)

	var (
		buttonHolder      = make(chan bool)
		backButtonPressed = false
	)

	// OK, now we process all events that come up
	for event := range fingerEvents {
		switch event {
		case listener.FINGER_UP:
			// The finger was just lifted from the sensor

			// This sends - if possible - an abort signal to the button holding goroutine below
			select {
			case buttonHolder <- false:
			default:
			}

			if backButtonPressed {
				backButtonPressed = false
			} else {
				if time.Since(lastTapTime) < 350*time.Millisecond {
					// This prevents pressing the power button twice quickly, which would open the default camera
					if time.Since(lastPowerTime) > 350*time.Millisecond {

						fmt.Println("Power")
						err = input.PressPowerButton(keyDevice)
						if err != nil {
							panic("pressing power button: " + err.Error())
						}
						lastPowerTime = time.Now()
						err := exec.Command("input", "keyevent", "KEYCODE_BACK").Run()
						if err != nil {
							panic("pressing back button: " + err.Error())
						}

					}
				} else {
					// Top left coordinates

					fmt.Println("Topleft")
					err = input.TouchUpDown(touchDevice, 35, 105)
					if err != nil {
						panic("running touch command: " + err.Error())
					}
				}
				lastTapTime = time.Now()
			}
		case listener.FINGER_DOWN:
			go func() {
				select {
				case <-buttonHolder:
					return
				case <-time.After(250 * time.Millisecond):

					fmt.Println("Back")
					backButtonPressed = true

					err := exec.Command("input", "keyevent", "KEYCODE_BACK").Run()
					if err != nil {
						panic("pressing back button: " + err.Error())
					}
					err = input.Vibrate(vibratorDevice, 50)
					if err != nil {
						panic("cannot vibrate: " + err.Error())
					}
				}
			}()

		}
	}
}
