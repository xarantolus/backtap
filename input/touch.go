package input

import (
	"encoding/binary"
	"io"
)

func RunTouch(output io.Writer, x, y uint32) (err error) {
	var touch = []InputEvent{
		{
			Type:  EV_ABS,
			Code:  ABS_MT_TRACKING_ID,
			Value: 0x0000e800, // TODO: Change this
		},
		{
			Type:  EV_KEY,
			Code:  BTN_TOUCH,
			Value: TOUCH_VALUE_DOWN,
		},
		{
			Type:  EV_KEY,
			Code:  BTN_TOOL_FINGER,
			Value: TOUCH_VALUE_DOWN,
		},
		{
			Type:  EV_ABS,
			Code:  ABS_MT_POSITION_X,
			Value: 0x00000071,
		},
		{
			Type:  EV_ABS,
			Code:  ABS_MT_POSITION_Y,
			Value: 0x000000a3,
		},
		{
			Type:  EV_ABS,
			Code:  ABS_MT_TOUCH_MAJOR,
			Value: 0x00000005,
		},
		eventSynReport,
		{
			Type:  EV_ABS,
			Code:  ABS_MT_TRACKING_ID,
			Value: 0xffffffff,
		},
		{
			Type:  EV_KEY,
			Code:  BTN_TOUCH,
			Value: TOUCH_VALUE_UP,
		},
		{
			Type:  EV_KEY,
			Code:  BTN_TOOL_FINGER,
			Value: TOUCH_VALUE_UP,
		},
		eventSynReport,
	}

	for _, ievent := range touch {
		err = binary.Write(output, binary.LittleEndian, ievent)
		if err != nil {
			return err
		}
	}

	return
}
