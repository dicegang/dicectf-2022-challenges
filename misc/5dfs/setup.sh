#!/bin/sh
./5dfswmtt &
sleep 0.1
rm /flag
exec setpriv --reuid 1000 --regid 1000 --init-groups bash
