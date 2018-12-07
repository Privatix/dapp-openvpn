#!/bin/sh

if [ -n "$1" ]
then
	server=$1
else
	echo "No arguments supplied (server)."
	exit 1
fi

# defines default internet interface
default=$(/sbin/route get default| grep interface| awk '{print $2}')

# creates rules
rules="nat on $default from $server/24 to any -> ($default)\nnat on $dev from $server/24 to any -> ($dev)"
rm -f ./config/nat_rules
echo $rules >> ./config/nat_rules

echo $default $dev $server

# enables ip forwarding
#/usr/sbin/sysctl -w net.inet.ip.forwarding=1
#/usr/sbin/sysctl -w net.inet.ip.fw.enable=1

#disables pfctl
#/sbin/pfctl -d
#sleep 1

#flushes all pfctl rules
#/sbin/pfctl -F all
#sleep 1

#starts pfctl and loads the rules from the nat-rules file
#/sbin/pfctl -f ./config/nat-rules -e



