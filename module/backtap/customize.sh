# Set to true if you do *NOT* want Magisk to mount
# any files for you. Most modules would NOT want
# to set this flag to true
SKIPMOUNT=false

# Set to true if you need to load system.prop
PROPFILE=false

# Set to true if you need post-fs-data script
POSTFSDATA=false

# Set to true if you need late_start service script
LATESTARTSERVICE=true

print_modname() {
  ui_print " "
  ui_print "       ********************************************"
  ui_print "       *           backtap by xarantolus          *"
  ui_print "       ********************************************"
  ui_print " "
}

set_permissions() {
  # The following is the default rule, DO NOT remove
  set_perm_recursive $MODPATH 0 0 0755 0644

  # Custom permissions
  set_perm $MODPATH/system/bin/backtap 0 2000 0755 u:object_r:system_file:s0
  set_perm $MODPATH/service.sh 0 2000 0755 u:object_r:system_file:s0
}
