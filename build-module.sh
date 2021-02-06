#!/bin/bash

set -e

# We just build the executable
bash ./build.sh

mkdir -p module/backtap/system/bin

# move it to our module
mv backtap module/backtap/system/bin/backtap

cd module/backtap

ZIP_NAME="MagiskModule-backtap.zip"

if [ "$1" == "-debug" ]
then
    echo "Building in debug mode"

    ZIP_NAME="MagiskModule-backtap-debug.zip"

    # Add file to enable debug mode (see service.sh)
    touch "DEBUG"
    # Read module props file before
    FILE_BEFORE=$(<module.prop)

    # exit_cleanup resets all of them on exit    
    exit_cleanup() {
        rm -f DEBUG
        echo "$FILE_BEFORE" > module.prop
    }
    trap exit_cleanup EXIT

    # Mark version as debug version
    sed -i -Ee 's/version=(.*)/version=\1-debug/g' module.prop
fi

# remove the old packed module, if possible
rm -f "../../$ZIP_NAME"

# and create the new one
zip -r "../../$ZIP_NAME" . 


