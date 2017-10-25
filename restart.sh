#!/bin/sh

pkill -9 ovs

cd /home/vyos/rvm
nohup ./ovs &

exit 0
