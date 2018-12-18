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

    # defines openvpn interface
    ip=${server%0}
    for interface in $(/sbin/ifconfig | grep 'utun\|inet.*-->' | sed -E 's/[[:space:]:].*//;/^$/d')
    do
        var=$(/sbin/ifconfig "$interface" | sed 1d | grep inet | grep "$ip")
        if [ ! -z "$var" ]
        then
            tun="$interface"
            break 2
        fi
    done

    if [ -z "$tun" ]
    then
        echo "openvpn interface not found."
	    exit 1
    fi


    # creates rules
    rm -f /usr/local/nat-rules

    nats="nat on $default from $server/24 to any -> ($default)\nnat on $tun from $server/24 to any -> ($tun)"
    echo "$nats" >> /usr/local/nat-rules

    ports="\npass in proto { tcp, udp } from any to any port $port"
    echo "$ports" >> /usr/local/nat-rules

    # enables ip forwarding
    /usr/sbin/sysctl -w net.inet.ip.forwarding=1

    #disables pfctl
    /sbin/pfctl -d
    sleep 1

    #flushes all pfctl rules
    /sbin/pfctl -F all
    sleep 1

    #starts pfctl and loads the rules from the nat-rules file
    /sbin/pfctl -f /usr/local/nat-rules -e
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
