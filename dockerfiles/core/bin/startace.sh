#!/bin/bash

docker run -d --privileged --name ace_rngd edge/rngd:1.0 # Needed to create entropy more quickly
ACE_CONSOLE=$(docker images | grep console-alpine | head -n 1 | awk '{print $3}')
export UUID=$(docker run -i --rm --entrypoint="" ${ACE_CONSOLE} uuidgen)
if grep -q localhost /etc/hosts; then echo "" > /dev/null; else echo "127.0.0.1       localhost" >> /etc/hosts; fi
docker-compose -p ace -f /opt/ace/compose/docker-compose.yml up -d
