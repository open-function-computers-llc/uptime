#! /bin/bash

echo "Killing any old processes"
killall ofc-uptime > /dev/null 2>&1

echo "Purging old binaries"
rm dist/ofc-uptime

echo "Building new binaries"
cd httpd
go build -o ../dist/ofc-uptime
cd ..

echo "Starting new built binary"
dist/ofc-uptime > log &
