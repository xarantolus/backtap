echo "Building"
./build.sh

echo "Pushing"
adb push -p "backtap" 'sdcard/'

echo "Install"
# adb shell "su -c \"killall backtap\" || true" || true

adb shell su -c "mv /sdcard/backtap /sbin"
adb shell su -c "chmod +x /sbin/backtap"
