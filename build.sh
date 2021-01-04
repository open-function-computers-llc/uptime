#! /bin/bash

# read in .env values
if [[ ! -f .env ]] ; then
    echo "You are missing the required .env file."
    exit 1
fi
source .env

export VERSION=`git rev-parse --short HEAD`

echo "Migrating the database"
cd migrations
goose mysql "$DB_USER:$DB_PASSWORD@tcp($DB_HOST:$DB_PORT)/$DB_NAME?parseTime=true" up
cd ..

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
