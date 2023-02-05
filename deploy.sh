echo "Building"
./build.sh

echo "Pushing"
adb push -p "backtap" "sdcard/"

echo "Install"

adb shell su -c "chmod +x /sdcard/backtap"

adb shell su -c "cp /sdcard/backtap /system/bin/backtap.dev"
adb shell su -c "chmod +x /system/bin/backtap.dev"
