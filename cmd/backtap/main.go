package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/xarantolus/backtap/input"
	"github.com/xarantolus/backtap/listener"
)

var debugMode bool

func main() {
	// Output debug timestamps with milliseconds
	log.SetFlags(log.Flags() | log.Lmicroseconds)
	flag.BoolVar(&debugMode, "debug", false, "Enable debug mode")
	flag.Parse()

	// we use this context to make sure all processes we start etc.
	// are closed/killed correctly in case something goes wrong
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	debug("starting listener")

	// Set up the stream we're listening to
	fingerEvents, err := listener.LogCat(ctx)
	if err != nil {
		panic("setting up logcat listener: " + err.Error())
	}

	// Open all device files we want to write to

	// Vibrator: When we write a number as string, the devices vibrates for that amount in milliseconds
	debug("opening vibrator device")
	vibratorDevice, err := os.OpenFile("/sys/devices/virtual/timed_output/vibrator/enable", os.O_WRONLY, os.ModeDevice)
	if err != nil {
		panic("cannot open vibrator device: " + err.Error())
	}
	defer vibratorDevice.Close()

	// Key: This input allows us to simulate a power button press
	debug("opening key device")
	keyDevice, err := os.OpenFile("/dev/input/event0", os.O_WRONLY, os.ModeDevice)
	if err != nil {
		panic("opening touch device input file: " + err.Error())
	}
	defer keyDevice.Close()

	// Touch: This input allows us to tap on the screen
	// This also works when the lockscreen is shown, so we should
	// be rather careful not to call emergency services accidentally
	debug("opening touchscreen device")
	touchDevice, err := os.OpenFile("/dev/input/event1", os.O_WRONLY, os.ModeDevice)
	if err != nil {
		panic("opening touch device input file: " + err.Error())
	}
	defer touchDevice.Close()

	debug("opening home button device")
	homeDevice, err := os.OpenFile("/dev/input/event2", os.O_WRONLY, os.ModeDevice)
	if err != nil {
		panic("opening home button device input file: " + err.Error())
	}
	defer homeDevice.Close()

	// The last time a certain action was performed
	var (
		lastTapTime   time.Time
		lastPowerTime time.Time
	)

	// This is for the "BACK" command. Since it needs to wait until the finger is on the sensor
	// for a few moments, its state management is a bit more complicated
	var (
		buttonAbort       = make(chan bool)
		backButtonPressed = false
		backButtonLock    sync.Mutex
	)

	debug("start processing events")

	// OK, now we process all events that come up
	for event := range fingerEvents {
		switch event {
		case listener.FINGER_UP:
			debug("FINGER_UP fired")

			// The finger was just lifted from the sensor

			// This sends - if it's running - an abort signal to the button holding goroutine below
			select {
			case buttonAbort <- false:
			default:
				{
				}
			}

			backButtonLock.Lock()
			if backButtonPressed {
				debug("Not clicking because we just ran the HOME command")
				backButtonPressed = false
				backButtonLock.Unlock()
			} else {
				backButtonLock.Unlock()

				if time.Since(lastTapTime) < 350*time.Millisecond {
					// This prevents pressing the power button twice quickly, which would open the default camera
					if time.Since(lastPowerTime) > 750*time.Millisecond {
						debug("Running SCREENOFF command")

						err = input.PressPowerButton(keyDevice)
						if err != nil {
							panic("pressing power button: " + err.Error())
						}
						debug("Finished SCREENOFF command")

						lastPowerTime = time.Now()
					} else {
						debug("Skipped SCREENOFF command as the screen is currently turning off")
					}
				} else {
					// Top left coordinates
					debug("Running TOUCH command")
					err = input.TouchUpDown(touchDevice, 50, 125)
					if err != nil {
						panic("running touch command: " + err.Error())
					}
					debug("Finished TOUCH command")
				}
				lastTapTime = time.Now()
			}
		case listener.FINGER_DOWN:
			debug("FINGER_DOWN fired")

			var started = make(chan bool)

			go func() {
				started <- true
				select {
				case <-buttonAbort:
					debug("aborted HOME command")
					return
				case <-time.After(250 * time.Millisecond):
					debug("Running HOME command")

					backButtonLock.Lock()
					backButtonPressed = true
					backButtonLock.Unlock()

					err := exec.Command("input", "keyevent", "KEYCODE_HOME").Run()
					if err != nil {
						panic("pressing home button: " + err.Error())
					}
					err = input.Vibrate(vibratorDevice, 50)
					if err != nil {
						panic("cannot vibrate: " + err.Error())
					}
					debug("Finished HOME command")
				}
			}()

			// Wait until the goroutine is actually running
			<-started
			close(started)
		}
	}
}

func debug(s ...interface{}) {
	if debugMode {
		log.Println(s...)
	}
}
