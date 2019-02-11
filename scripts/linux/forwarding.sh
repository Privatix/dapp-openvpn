#!/bin/sh

frwd=$(/sbin/sysctl -n net.ipv4.ip_forward)
if [ "$frwd" != "1" ]
then
    DIRECTORY=`dirname $0`

    cd ${DIRECTORY}
    cd ../../../
    
    # enables ip forwarding in pre-start.sh
    echo "\nsudo /sbin/sysctl -w net.ipv4.ip_forward=1\n" >> ./dappctrl/pre-start.sh

    # disable ip forwarind in post-stop.sh
    echo "\nsudo /sbin/sysctl -w net.ipv4.ip_forward=0\n" >> ./dappctrl/post-stop.sh
fi
