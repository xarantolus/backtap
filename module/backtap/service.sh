#!/system/bin/sh

# This script runs on boot. Here we can start our service
MODDIR=${0%/*}

# Log debug message 
if [ -f "$MODDIR/DEBUG" ]; then
  echo "Waiting for user to unlock/decrypt phone" > /cache/backtap.log
fi

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

# Make backtap executable in case it isn't
BTPATH="/system/bin/backtap"
chmod +x "$BTPATH"


if [ -f "$MODDIR/DEBUG" ]; then
  # Start in debug mode and log to cache directory
  nohup "$BTPATH" -debug >>/cache/backtap.log 2>&1 &
else
  # Start in normal mode
  nohup "$BTPATH" &
fi

