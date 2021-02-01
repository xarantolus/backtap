#!/system/bin/sh

# This script runs on boot. Here we can start our service
MODDIR=${0%/*}

echo "Waiting for user to unlock/decrypt phone" > /cache/backtap.log

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

BTPATH="/system/bin/backtap"

chmod +x "$BTPATH"

nohup "$BTPATH" -debug >>/cache/backtap.log 2>&1 &
