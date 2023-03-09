#!/bin/bash

docker-compose -p ace -f /opt/ace/compose/docker-compose.yml down -v
docker stop ace_rngd && docker rm ace_rngd

for x in $(ls /mnt/); do
    umount /mnt/$x >/dev/null
done
