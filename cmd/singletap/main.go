package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"syscall"
)

// We need the input_event struct from https://android.googlesource.com/platform/system/core/+/froyo-release/toolbox/sendevent.c

/*
struct input_event {
	struct timeval time;
	__u16 type;
	__u16 code;
	__s32 value;
};
*/

/*
add device 7: /dev/input/event1
  name:     "synaptics_dsx"
  events:
    KEY (0001): KEY_WAKEUP            BTN_TOOL_FINGER       BTN_TOUCH
    ABS (0003): ABS_X                 : value 0, min 0, max 1079, fuzz 0, flat 0, resolution 0
                ABS_Y                 : value 0, min 0, max 2159, fuzz 0, flat 0, resolution 0
                ABS_MT_SLOT           : value 0, min 0, max 9, fuzz 0, flat 0, resolution 0
                ABS_MT_TOUCH_MAJOR    : value 0, min 0, max 255, fuzz 0, flat 0, resolution 0
                ABS_MT_TOUCH_MINOR    : value 0, min 0, max 255, fuzz 0, flat 0, resolution 0
                ABS_MT_POSITION_X     : value 0, min 0, max 1079, fuzz 0, flat 0, resolution 0
                ABS_MT_POSITION_Y     : value 0, min 0, max 2159, fuzz 0, flat 0, resolution 0
                ABS_MT_TRACKING_ID    : value 0, min 0, max 65535, fuzz 0, flat 0, resolution 0
  input props:
	INPUT_PROP_DIRECT

add device 7: /dev/input/event1
  name:     "synaptics_dsx"
  events:
    KEY (0001): 008f  0145  014a
    ABS (0003): 0000  : value 0, min 0, max 1079, fuzz 0, flat 0, resolution 0
                0001  : value 0, min 0, max 2159, fuzz 0, flat 0, resolution 0
                002f  : value 0, min 0, max 9, fuzz 0, flat 0, resolution 0
                0030  : value 0, min 0, max 255, fuzz 0, flat 0, resolution 0
                0031  : value 0, min 0, max 255, fuzz 0, flat 0, resolution 0
                0035  : value 0, min 0, max 1079, fuzz 0, flat 0, resolution 0
                0036  : value 0, min 0, max 2159, fuzz 0, flat 0, resolution 0
                0039  : value 0, min 0, max 65535, fuzz 0, flat 0, resolution 0
  input props:
	INPUT_PROP_DIRECT
*/

/*
Putting finger in top left corner:

/dev/input/event1: 0003 0039 000000e7 -
/dev/input/event1: 0001 014a 00000001 -
/dev/input/event1: 0001 0145 00000001 -
/dev/input/event1: 0003 0035 00000078 -
/dev/input/event1: 0003 0036 000000e2 -
/dev/input/event1: 0003 0030 00000006 -
/dev/input/event1: 0000 0000 00000000 -
/dev/input/event1: 0003 0039 ffffffff -
/dev/input/event1: 0001 014a 00000000
/dev/input/event1: 0001 0145 00000000
/dev/input/event1: 0000 0000 00000000

Same thing:

/dev/input/event1: EV_ABS       ABS_MT_TRACKING_ID   000000e8 - *
/dev/input/event1: EV_KEY       BTN_TOUCH            DOWN     - *
/dev/input/event1: EV_KEY       BTN_TOOL_FINGER      DOWN     - *
/dev/input/event1: EV_ABS       ABS_MT_POSITION_X    0000007a -*
/dev/input/event1: EV_ABS       ABS_MT_POSITION_Y    000000a3 -*
/dev/input/event1: EV_ABS       ABS_MT_TOUCH_MAJOR   00000005 -*
/dev/input/event1: EV_SYN       SYN_REPORT           00000000 -*
/dev/input/event1: EV_ABS       ABS_MT_TOUCH_MAJOR   00000006 ?
/dev/input/event1: EV_SYN       SYN_REPORT           00000000 ?
/dev/input/event1: EV_ABS       ABS_MT_TRACKING_ID   ffffffff -*
/dev/input/event1: EV_KEY       BTN_TOUCH            UP
/dev/input/event1: EV_KEY       BTN_TOOL_FINGER      UP
/dev/input/event1: EV_SYN       SYN_REPORT           00000000
*/

// that's it. That's the struct
type InputEvent struct {
	Time  syscall.Timeval
	Type  EventType
	Code  EventCode
	Value uint32
}

// Here are some constants I infered from the output above
type EventType uint16

const (
	EV_ABS EventType = 0x0003
	EV_KEY EventType = 0x0001
	EV_SYN EventType = 0x0000
)

type EventCode uint16

const (
	ABS_MT_TRACKING_ID EventCode = 0x0039
	BTN_TOUCH          EventCode = 0x014a
	BTN_TOOL_FINGER    EventCode = 0x0145
	ABS_MT_POSITION_X  EventCode = 0x0035
	ABS_MT_POSITION_Y  EventCode = 0x0036
	ABS_MT_TOUCH_MAJOR EventCode = 0x0030
	SYN_REPORT         EventCode = 0x0000
)

const (
	TOUCH_VALUE_DOWN = 0x00000001
	TOUCH_VALUE_UP   = 0x00000000
)

var (
	eventSynReport = InputEvent{
		Type:  EV_SYN,
		Code:  SYN_REPORT,
		Value: 0x00000000,
	}
)

var (
	touch = []InputEvent{
		{
			Type:  EV_ABS,
			Code:  ABS_MT_TRACKING_ID,
			Value: 0x0000e800, // TODO: Change this - touch tracking ID
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
		// Top right corner
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
)

func main() {
	f, err := os.OpenFile("/dev/input/event1", os.O_WRONLY, os.ModeDevice)
	if err != nil {
		panic("opening input file: " + err.Error())
	}
	defer f.Close()

	for _, ievent := range touch {
		err = binary.Write(f, binary.LittleEndian, ievent)
		if err != nil {
			panic("writing to input event: " + err.Error())
		}
	}
	fmt.Println("Done")
}
