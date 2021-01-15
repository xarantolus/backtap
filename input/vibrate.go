package input

import (
	"io"
	"strconv"
)

func Vibrate(device io.Writer, ms int) (err error) {
	_, err = device.Write([]byte(strconv.Itoa(ms)))
	return
}
