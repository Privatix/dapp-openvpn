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

    # defines default internet interface
    default=$(/sbin/route | grep  default | awk '{print $8}')

    # defines openvpn interface
    ip=${server%0}
    for interface in $(/sbin/ifconfig | grep 'tun\|inet.*-->' | sed -E 's/[[:space:]:].*//;/^$/d')
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

    /bin/sed -i 's/^dns=dnsmasq/#&/' /etc/NetworkManager/NetworkManager.conf && /usr/sbin/service network-manager restart

    # creates rules
    /sbin/iptables -t nat -A POSTROUTING -s $server/24 -o $default -j MASQUERADE
elif [ "$status" = "off" ]
then
    if [ -n "$2" ]
    then
	    server=$2
    else
	    echo "No arguments supplied (server)."
	    exit 1
    fi

    # defines default internet interface
    default=$(/sbin/route | grep  default | awk '{print $8}')

    /sbin/iptables -t nat -D POSTROUTING -s $server/24 -o $default -j MASQUERADE
fi
