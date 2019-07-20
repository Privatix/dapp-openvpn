#!/bin/bash
# 
# Parses DHCP options from openvpn to update resolv.conf
# To use set as 'up' and 'down' script in your openvpn *.conf:
# up /etc/openvpn/update-resolv-conf
# down /etc/openvpn/update-resolv-conf
#
# Used snippets of resolvconf script by Thomas Hood and Chris Hanson.
# Licensed under the GNU GPL.  See /usr/share/common-licenses/GPL. 
# 
# Example envs set from openvpn:
#
#     foreign_option_1='dhcp-option DNS 193.43.27.132'
#     foreign_option_2='dhcp-option DNS 193.43.27.133'
#     foreign_option_3='dhcp-option DOMAIN be.bnc.ch'
#

[ -x /etc/resolv.conf ] || exit 0
[ "$script_type" ] || exit 0
[ "$dev" ] || exit 0

split_into_parts()
{
	part1="$1"
	part2="$2"
	part3="$3"
}

NMSRVRS=""
for optionvarname in ${!foreign_option_*} ; do
         option="${!optionvarname}"
         echo "UP script: $option"
         split_into_parts $option
         if [ "$part1" = "dhcp-option" ] ; then
         	if [ "$part2" = "DNS" ] ; then
                	NMSRVRS="${NMSRVRS:+$NMSRVRS }$part3"
                fi
         fi
done

case "$script_type" in
  up)
	for NS in $NMSRVRS ; do
        	R="${R}nameserver $NS
"
	done
	echo -n "$R" >> /etc/resolv.conf
	;;
  down)
        for NS in $NMSRVRS ; do
                echo "$(sed '0,/'$NS'/{/'$NS'/d;}' /etc/resolv.conf)" >  /etc/resolv.conf		
        done
	;;
esac

