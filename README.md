# backtap
Phones have become larger in recent years. While this has in general been an improvement, some things have been left behind. When was the last time you wanted to tap a button that was out of the reach of your thumb? And when did you last drop your phone while trying to reach that button? Also, why is the fingerprint reader so useless after the phone has been unlocked?

What if we could combine all these problems into one solution? That's what `backtap` does.

When the phone is locked, the fingerprint reader serves its normal purpose of unlocking the phone. When the phone is unlocked, it gains ✨*special abilities*✨.

There are three different commands:
* Tapping the fingerprint reader one time is the same as clicking in the top left corner of the phone. This is where menu, back or home buttons are positioned. That way, instead of reaching the top of the phone, you just tap the fingerprint reader.
* Tapping the fingerprint reader twice (*tap tap*) does the same as pressing the off button of your phone, it just turns off the screen. This is especially convenient when the off button is located at the side of the phone, which makes it hard to reach.
* Holding the fingerprint sensor a bit longer opens the "recent apps" menu. That makes switching apps easier, especially if your phone doesn't have any buttons at the front.

### Supported phones
Whether or not this works on your phone might not only depend on the model used, but also the operating system. Different variants of Android might not behave in the way this module expects them to.

The only system I tested this on is [LineageOS 17.1](https://lineageos.org/) for the [Xiaomi Mi Mix 2](https://wiki.lineageos.org/devices/chiron) in combination with [Magisk 21.4](https://github.com/topjohnwu/Magisk).

The program is very device specific and will very likely not work at all if you don't have the *exact* same configuration. It should however be possible to [adapt the program](#adapting-to-other-phones) to run on your phone.

This program *might* work on your phone if all of these conditions are met:
* You have Magisk installed
* The command `echo -n 200 > /sys/devices/virtual/timed_output/vibrator/enable` in a root shell vibrates your phone
* You have adb logs enabled ("Log buffer size" in developer options should *not* be set to "Off")
* The `singletap` executable (that can be built from [`cmd/singletap/main.go`](cmd/singletap/main.go)) correctly taps the top left of your screen 
* The output of `getevent -pl` (again, in a root shell on your phone) looks something like this (`/dev/input/event0` (power button) and `/dev/input/event1` (display) are particularly important)

<details>
<summary>Expected output</summary>
<pre>
chiron:/ # getevent -pl
add device 1: /dev/input/event6
  name:     "msm8998-tasha-snd-card Button Jack"
  events:
    KEY (0001): KEY_VOLUMEDOWN        KEY_VOLUMEUP          KEY_MEDIA             BTN_3
                BTN_4                 BTN_5
  input props:
    INPUT_PROP_ACCELEROMETER
add device 2: /dev/input/event5
  name:     "msm8998-tasha-snd-card Headset Jack"
  events:
    SW  (0005): SW_HEADPHONE_INSERT   SW_MICROPHONE_INSERT  SW_LINEOUT_INSERT     SW_JACK_PHYSICAL_INS
                SW_PEN_INSERTED       0010                  0011                  0012
  input props:
    <none>
add device 3: /dev/input/event4
  name:     "uinput-fpc"
  events:
    KEY (0001): KEY_KPENTER           KEY_UP                KEY_LEFT              KEY_RIGHT
                KEY_DOWN              BTN_GAMEPAD           BTN_EAST              BTN_C
                BTN_NORTH             BTN_WEST
  input props:
    <none>
add device 4: /dev/input/event0
  name:     "qpnp_pon"
  events:
    KEY (0001): KEY_VOLUMEDOWN        KEY_POWER
  input props:
    <none>
add device 5: /dev/input/event3
  name:     "gpio-keys"
  events:
    KEY (0001): KEY_VOLUMEUP
    SW  (0005): SW_LID
  input props:
    <none>
add device 6: /dev/input/event2
  name:     "uinput-goodix"
  events:
    KEY (0001): KEY_HOME
  input props:
    <none>
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
</pre>
</details>


### How it works
I have tried several variations of this idea for quite some time. This program is basically just the next iteration, but likely the final one as this is as low-level as it gets.

#### First idea
My first implementation was done using [Tasker](https://play.google.com/store/apps/details?id=net.dinglisch.android.taskerm), an automation app. I only implemented the top left click as anything more would have become very complicated with visual programming.
The concept was simple: Whenever a certain line appeared in the system log (indicating that a finger was put on/off the sensor), it would run the command `input tap x y`, which taps the specified coordinates of the screen.

The problem with this was that this method is *noticeably* slow. As in I clicked the sensor, waited for some time and only then the action happend. It took about 500-1000ms, which was just too slow to be usable.

#### Speeding up
In my search for a faster solution I took a look at what the input command actually does: it is a shell script starting a whole Java program that writes a quite small signal to the display. You can cut out two middlemen if done right.

So I took a look on *what* was written. If that program can write the required signal, so could a faster one. That's how this project was born.

I basically used the Android command-line tool `getevent` to see what was written and noticed that I would have to write certain bits to certain files to make it work.

To be more specific, one has to write an [`input_event`](https://android.googlesource.com/platform/system/core/+/froyo-release/toolbox/sendevent.c#13) which looks like this
```c
struct input_event {
	struct timeval time;
	__u16 type;
	__u16 code;
	__s32 value;
};
```
to an input device file that is located at `/dev/input/`.

After quite some time of reading the `getevent` output, I implemented the [same struct in Go](input/event.go#L90). This data can be written to a device file, which lets the device interpret the signal.

I then translated the `getevent` output of a click into versions of that struct, which allowed me to simulate clicks on the display. You can see this functionality in the [`singletap`](cmd/singletap/main.go) program of this repo.

And the speed difference is quite nice. Instead of waiting for more than half a second, the new click method was so fast that I didn't even notice the delay.

#### Combining
I then combined this method of clicking and a new method to read the system log for fingerprint sensor events into one program, which is located at [`cmd/backtap/main.go`](cmd/backtap/main.go).

If the program is running, it will detect clicks on the fingerprint sensor. This detection might not seem optimal (it basically just scans the `adb logcat` output for certain lines), but it is surprisingly fast. Like, *very* fast.

Whenever an event happens, the program interprets it and either starts a *very fast* click, a *quite fast* screen off or a *quite slow* "recent" button command (this last one still uses the `input` command). And that's exactly what I wanted.

### Building
Since this is a [Go](https://golang.org/) program, compiling the code should be quite easy (I promise!). You do need to know the architecture your phone uses, but it's likely `arm64`.

If that's the case, you can run...
* [`build.sh`](build.sh) to build the program
* [`deploy.sh`](deploy.sh) to build the program and push it to your phone (if it has an `/sbin` directory, you can now use it from the command line)
* [`build-module.sh`](build-module.sh) to package this program in a Magisk Module you can install from Magisk Manager. It will run `backtap` on boot. To build a debug version that creates a log file at `/cache/backtap.log`, you can pass `-debug` to this script as first parameter.

If not, you have to look into what the `GOARCH` environment variable does and then edit `build.sh` to use the correct value for your phone.

### Adapting to other phones
Adapting this program to run on other phone models and configurations should be possible. You would have to change the `logcat` line that is detected and the way these commands are executed. It is likely that the names of input devices differ. The way the touchscreen tap works should be the same, as that's based on on a protocol that seems to have been used for quite a long time by many different companies (since it's part of Linux).

### [License](LICENSE)
This is free as in freedom software. Do whatever you like with it.
