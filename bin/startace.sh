#!/bin/bash

docker run -d --privileged --name ace_rngd edge/rngd:1.0 # Needed to create entropy more quickly
export UUID=$(docker run -i --rm --entrypoint="" edge/console-alpine:1.0 uuidgen)
if grep -q localhost /etc/hosts; then echo "" > /dev/null; else echo "127.0.0.1       localhost" >> /etc/hosts; fi
docker-compose -p ace -f /opt/ace/compose/docker-compose.yml up -d
