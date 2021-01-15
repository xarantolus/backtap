package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/xarantolus/backtap/input"
)

const (
	logcatCommand = "logcat"
)

// Lines we're looking for:

// Finger is on the sensor:
//     report_input_event - Reporting event type: 1, code: 96, value:1
// Finger is off the sensor after being on it:
//     report_input_event - Reporting event type: 1, code: 96, value:0
// Afterwards:
// nav_loop waiting for finger down

func main() {
	var (
		startTime = time.Now()
		pauseMode = true
	)

	fmt.Println("Program running")
	ctx, cancel := context.WithCancel(context.Background())
	// makes sure the logcat command is killed afterwards
	defer cancel()

	// --format=raw makes sure we don't get unnecessary time info
	var watcher = exec.CommandContext(ctx, logcatCommand, "--format=raw")

	stdout, err := watcher.StdoutPipe()
	if err != nil {
		panic("cannot connect stdout pipe: " + err.Error())
	}

	err = watcher.Start()
	if err != nil {
		panic("starting logcat: " + err.Error())
	}

	vibratorDevice, err := os.OpenFile("/sys/devices/virtual/timed_output/vibrator/enable", os.O_WRONLY, os.ModeDevice)
	if err != nil {
		panic("cannot open vibrator device: " + err.Error())
	}
	defer vibratorDevice.Close()

	keyDevice, err := os.OpenFile("/dev/input/event0", os.O_WRONLY, os.ModeDevice)
	if err != nil {
		panic("opening touch device input file: " + err.Error())
	}
	defer keyDevice.Close()

	touchDevice, err := os.OpenFile("/dev/input/event1", os.O_WRONLY, os.ModeDevice)
	if err != nil {
		panic("opening touch device input file: " + err.Error())
	}
	defer touchDevice.Close()

	scan := bufio.NewScanner(stdout)

	var (
		lastTapTime   time.Time
		lastPowerTime time.Time
	)

	var (
		buttonHolder      = make(chan bool)
		backButtonPressed = false
	)
	for scan.Scan() {
		// The first few seconds of running should not process anything
		if pauseMode {
			if time.Since(startTime) > 5*time.Second {
				pauseMode = false
				fmt.Println("Now accepting commands")
			}
			continue
		}

		text := scan.Text()

		if strings.HasPrefix(text, "report_input_event") {
			lastChar := text[len(text)-1]

			if lastChar == '0' {
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
			} else if lastChar == '1' {
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

}
