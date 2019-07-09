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

    # Add , block from VPN to LAN

    block_tolan="rfc1918 = \"{ 192.168.0.0/16, 172.16.0.0/12, 10.0.0.0/8 }\" \nvpnnet = \"{ 10.217.3.0/24 }\" \nblock in log quick from \$vpnnet to \$rfc1918"
    echo "${block_tolan}" >> /usr/local/nat-rules

    frwd=$(/usr/sbin/sysctl -n net.inet.ip.forwarding)
    if [ "$frwd" != "1" ]
    then
        # enables ip forwarding
        /usr/sbin/sysctl -w net.inet.ip.forwarding=1
    fi

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
    if [ -n "$2" ]
    then
	    forwarding=$2
    else
	    echo "No arguments supplied (forwarding)."
	    exit 1
    fi

    if [ "$forwarding" = "0" ]
    then
        # disables ip forwarding
        /usr/sbin/sysctl -w net.inet.ip.forwarding=0
    fi

    #disables pfctl
    /sbin/pfctl -d
    sleep 1

    #flushes all pfctl rules
    /sbin/pfctl -F all
    sleep 1

    #starts pfctl and loads the default rules
    /sbin/pfctl -f /etc/pf.conf -e
fi
