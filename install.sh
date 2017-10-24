#!/bin/bash

HOME_DIR=/home/vyos/rvm
VYATTA_BOOTFILE=/opt/vyatta/etc/config/scripts/vyatta-postconfig-bootup.script

cat $VYATTA_BOOTFILE | grep ovsboot > /dev/null
RET=$?
if [ $RET = 0 ]; then
	exit 0
fi

cat >>$VYATTA_BOOTFILE <<EOF
#!/bin/bash

chmod +x $HOME_DIR/ovs
chmod +x $HOME_DIR/ovsboot
chown vyos:users $HOME_DIR/ovs

$HOME_DIR/ovsboot >/home/vyos/rvm/ovsboot.log 2>&1 &
exit 0

EOF

echo "successfully write ovs bootinfo to the system script"
