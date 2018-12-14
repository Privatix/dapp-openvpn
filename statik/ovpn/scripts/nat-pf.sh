#!/bin/sh

if [ -n "$1" ]
then
	status=$1
else
	echo "No arguments supplied (status)."
	exit 1
fi

if [ "$status" = "on" ]
then
    if [ -n "$2" ]
    then
	    server=$2
    else
	    echo "No arguments supplied (server)."
	    exit 1
    fi

    if [ -n "$3" ]
    then
	    port=$3
    else
	    echo "No arguments supplied (port)."
	    exit 1
    fi

    # defines default internet interface
    default=$(/sbin/route get default| grep interface| awk '{print $2}')

    # creates rules
    rm -f ./config/nat-rules

    nats="nat on $default from $server/24 to any -> ($default)\nnat on $dev from $server/24 to any -> ($dev)"
    echo "$nats" >> ./config/nat-rules

    ports="\npass in proto { tcp, udp } from any to any port $port"
    echo "$ports" >> ./config/nat-rules

    echo $default $dev $server $port

    # enables ip forwarding
    /usr/sbin/sysctl -w net.inet.ip.forwarding=1

    #disables pfctl
    /sbin/pfctl -d
    sleep 1

    #flushes all pfctl rules
    /sbin/pfctl -F all
    sleep 1

    #starts pfctl and loads the rules from the nat-rules file
    /sbin/pfctl -f ./config/nat-rules -e
elif [ "$status" = "off" ]
then
    # disables ip forwarding
    /usr/sbin/sysctl -w net.inet.ip.forwarding=0

    #disables pfctl
    /sbin/pfctl -d
    sleep 1

    #flushes all pfctl rules
    /sbin/pfctl -F all
    sleep 1

    #starts pfctl and loads the default rules
    /sbin/pfctl -f /etc/pf.conf -e
fi
