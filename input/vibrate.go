package input

import (
	"os"
	"path/filepath"
	"strconv"
)

// Vibrate vibrates the device (duh). The vibratorPath is **not** a device file,
// but rather the path to the device directory.
// It is probably /sys/class/leds/vibrator/, but it could depend on the device.
// Also see https://stackoverflow.com/a/62678837
func Vibrate(vibratorPath string, milliseconds int) (err error) {
	// At first, we write the duration (as string) to the file named "duration",
	// then we activate it by writing "1" to the file named "activate"

	// I do not know who thought that it is a good idea to have the
	// vibration device under LEDS, but it sure as hell is confusing

	err = writeInt(filepath.Join(vibratorPath, "duration"), milliseconds)
	if err != nil {
		return
	}

	return writeInt(filepath.Join(vibratorPath, "activate"), 1)
}

func writeInt(path string, i int) (err error) {
	f, err := os.OpenFile(path, os.O_WRONLY, os.ModeDevice)
	if err != nil {
		return
	}

	_, err = f.Write([]byte(strconv.Itoa(i)))
	if err != nil {
		_ = f.Close()
		return
	}

	return f.Close()
}
