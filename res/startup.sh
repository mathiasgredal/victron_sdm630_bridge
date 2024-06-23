#!/bin/bash
while true; do
        /data/bridge/mbmd run  -a /dev/ttyUSB0 -d sdm:1 -b 38400 --comset 8E1 --mqtt-broker='127.0.0.1:1883' --mqtt-topic='stromzaehler' &
        /data/bridge/sdm630-bridge-arm
        sleep 1
done