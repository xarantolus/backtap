package listener

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"log"
	"os/exec"
	"time"
)

type FingerEvent uint8

const (
	FINGER_DOWN FingerEvent = 1
	FINGER_UP   FingerEvent = 0
)

// LogCat watches the logcat stream for any fingerprint events
func LogCat(ctx context.Context) (events chan FingerEvent, err error) {
	// Maximum amount of events we want to buffer
	events = make(chan FingerEvent, 1)

	// --format=raw makes sure we don't get unnecessary time info
	// -T 1 makes sure we only get the latest lines, ignoring those since the boot
	// -e makes sure that logcat already filters most lines, which should mean that we need to process less

	var lastLogcatStart time.Time

	var startLogcat = func() (watcher *exec.Cmd, stdout io.ReadCloser, err error) {
		watcher = exec.CommandContext(ctx, "logcat", "--format=raw", "-T", "1", "-e", "report_input_event")

		// Logs can be read from stdout
		stdout, err = watcher.StdoutPipe()
		if err != nil {
			return
		}

		lastLogcatStart = time.Now()

		// Start the command
		err = watcher.Start()

		return
	}

	watcher, stdout, err := startLogcat()
	if err != nil {
		return
	}

	// Everything that comes now should be processed
	scan := bufio.NewScanner(stdout)

	// Lines we're looking for:
	//
	// Finger is on the sensor:
	//     report_input_event - Reporting event type: 1, code: 96, value:1
	// Finger is off the sensor after being on it:
	//     report_input_event - Reporting event type: 1, code: 96, value:0
	// Afterwards:
	// nav_loop waiting for finger down

	// Now we just scan in the background until the stream breaks off
	go func() {
		defer close(events)

		for {
			// Scan all lines in the background
			for scan.Scan() {
				text := scan.Bytes()

				if bytes.HasPrefix(text, []byte("report_input_event")) {
					lastChar := text[len(text)-1]

					// Process this last character after "value:"
					if lastChar == '0' {
						events <- FINGER_UP
					} else if lastChar == '1' {
						events <- FINGER_DOWN
					}
				}
			}

			// Panicing here instead of moving it to main makes the logic a bit easier
			err := scan.Err()
			if err != nil {
				panic("streaming logcat output: " + err.Error())
			}

			// Make sure resources are cleaned up
			_ = watcher.Wait()

			// Flush log buffers, this can help get logcat out of this state where we would have to restart often
			_ = exec.CommandContext(ctx, "logcat", "--clear").Run()

			// OK, not trying again
			if time.Since(lastLogcatStart) < 5*time.Second {
				panic("logcat crashed seconds after restarting it, exiting")
			}

			// Close everything and try again
			_ = stdout.Close()

			log.Println("logcat: restarting watcher")

			// Set up a new watcher
			watcher, stdout, err = startLogcat()
			if err != nil {
				panic("restarting logcat: " + err.Error())
			}
			scan = bufio.NewScanner(stdout)
		}
	}()

	return
}
