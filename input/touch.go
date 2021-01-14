package input

import (
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
			Value: DOWN,
		},
		{
			Type:  EV_KEY,
			Code:  BTN_TOOL_FINGER,
			Value: DOWN,
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
			Value: UP,
		},
		{
			Type:  EV_KEY,
			Code:  BTN_TOOL_FINGER,
			Value: UP,
		},
		eventSynReport,
	}

	return runEvents(output, touch)
}

func TouchDown(output io.Writer, x, y uint32) error {
	return runEvents(output, []InputEvent{
		{
			Type:  EV_ABS,
			Code:  ABS_MT_TRACKING_ID,
			Value: 0x0000e800, // TODO: Change this
		},
		{
			Type:  EV_KEY,
			Code:  BTN_TOUCH,
			Value: DOWN,
		},
		{
			Type:  EV_KEY,
			Code:  BTN_TOOL_FINGER,
			Value: DOWN,
		},
		{
			Type:  EV_ABS,
			Code:  ABS_MT_POSITION_X,
			Value: x,
		},
		{
			Type:  EV_ABS,
			Code:  ABS_MT_POSITION_Y,
			Value: y,
		},
		{
			Type:  EV_ABS,
			Code:  ABS_MT_TOUCH_MAJOR,
			Value: 0x00000005,
		},
		eventSynReport,
	})
}

func TouchUp(output io.Writer) error {
	return runEvents(output, []InputEvent{
		{
			Type:  EV_ABS,
			Code:  ABS_MT_TRACKING_ID,
			Value: 0xffffffff,
		},
		{
			Type:  EV_KEY,
			Code:  BTN_TOUCH,
			Value: UP,
		},
		{
			Type:  EV_KEY,
			Code:  BTN_TOOL_FINGER,
			Value: UP,
		},
		eventSynReport,
	})
}
