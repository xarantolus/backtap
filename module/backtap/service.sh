#!/system/bin/sh

# This script runs on boot. Here we can start our service
MODDIR=${0%/*}

# Wait until boot & decryption finished
until [ -d /sdcard/Download ]
do
  sleep 5
done
pgrep zygote > /dev/null && {
  until [ .$(getprop sys.boot_completed) = .1 ]
  do
    sleep 5
  done
}

chmod +x /system/bin/backtap

# Run the service in background!
backtap -debug >> /sdcard/backtap.log 2>&1 &


