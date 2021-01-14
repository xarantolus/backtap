package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
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

	shell := exec.CommandContext(ctx, "sh", "-li")

	shellIn, err := shell.StdinPipe()
	if err != nil {
		panic("connecting shell stdin pipe: " + err.Error())
	}
	defer shellIn.Close()

	err = shell.Start()
	if err != nil {
		panic("starting shell: " + err.Error())
	}

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

	touchDevice, err := os.OpenFile("/dev/input/event1", os.O_WRONLY, os.ModeDevice)
	if err != nil {
		panic("opening touch device input file: " + err.Error())
	}
	defer touchDevice.Close()

	scan := bufio.NewScanner(stdout)

	var lastTapTime time.Time

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

		if strings.HasPrefix(text, "report_input_event") && strings.HasSuffix(text, "0") {
			if time.Since(lastTapTime) < 750*time.Millisecond {
				err = writeShell(shellIn, "input keyevent 26")
				if err != nil {
					panic("turning off screen: " + err.Error())
				}
			} else {

				err = input.RunTouch(touchDevice, 35, 105)
				if err != nil {
					panic("running touch command: " + err.Error())
				}
			}
			lastTapTime = time.Now()
		}
	}

}

func writeShell(shell io.Writer, cmd string) (err error) {
	_, err = fmt.Fprintln(shell, cmd)
	return
}
