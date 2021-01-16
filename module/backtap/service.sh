# This script runs on boot. Here we can start our service
MODDIR=${0%/*}

# make sure it's executable
chmod +x $MODDIR/common/backtap

# Run the service!
$MODDIR/common/backtap
