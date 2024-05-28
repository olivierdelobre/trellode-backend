#!/bin/bash

docker-compose up -d --build
echo "Checking for local server to be started"

CPT=0
IS_RUNNING=0
while true
do
	sleep 3

	echo "checking...$CPT"

	CPT=$((CPT+1))
	if [ "$CPT" == 200 ]; then
    break
  fi

  STARTED=0
 
  RC=`nc -z localhost 8080;echo $?`

  if [[ "$RC" == "0" ]]; then
    STARTED=1
  fi

  if [ "$STARTED" == 1 ]; then
    IS_RUNNING=1
    break
  fi
done

if [ "$IS_RUNNING" == 1 ]; then
  echo "All up and running!"
  go test -v internal/api/api_test.go
  RC=`echo $?`
  exit $RC
else
  echo "Failed to run tests, local server could not start"
  exit 1
fi