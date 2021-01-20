set -e

# We just build the executable
./build.sh

mkdir -p module/backtap/system/bin

# move it to our module
mv backtap module/backtap/system/bin/backtap

cd module/backtap

# remove the old packed module, if possible
rm -f ../../MagiskModule-backtap.zip 

# and create the new one
zip -r ../../MagiskModule-backtap.zip .
