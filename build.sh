#! /bin/bash

echo "Killing any old processes"
killall uptime-osx > /dev/null 2>&1
killall uptime-linux > /dev/null 2>&1

echo "Purging old binaries"
rm dist/uptime*

echo "Building new binaries"
cd httpd
GOOS=darwin GOARCH=amd64 go build -o ../dist/uptime-osx
GOOS=linux GOARCH=amd64 go build -o ../dist/uptime-linux
cd ..

echo "Starting new built binary"
if [[ "$OSTYPE" == "linux-gnu" ]]; then
    dist/uptime-linux &
elif [[ "$OSTYPE" == "darwin"* ]]; then
    dist/uptime-osx &
else
    echo "i have no idea what's going on"
fi
