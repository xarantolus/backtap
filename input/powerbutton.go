package input

import (
	"encoding/binary"
	"io"
	"time"
)

/*
Press power button:

/dev/input/event0: 0001 0074 00000001
/dev/input/event0: 0000 0000 00000000
/dev/input/event0: 0001 0074 00000000
/dev/input/event0: 0000 0000 00000000

*/

const (
	POWER_BUTTON EventCode = 0x0074
)

func PressPowerButton(output io.Writer) (err error) {
	err = binary.Write(output, binary.LittleEndian, InputEvent{
		Type:  EV_KEY,
		Code:  POWER_BUTTON,
		Value: DOWN,
	})
	if err != nil {
		return err
	}
	err = binary.Write(output, binary.LittleEndian, eventSynReport)
	if err != nil {
		return err
	}
	time.Sleep(25 * time.Millisecond)
	err = binary.Write(output, binary.LittleEndian, InputEvent{
		Type:  EV_KEY,
		Code:  POWER_BUTTON,
		Value: UP,
	})
	if err != nil {
		return err
	}
	err = binary.Write(output, binary.LittleEndian, eventSynReport)

	return err
}
