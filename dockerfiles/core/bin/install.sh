#!/bin/bash

# we want to have some checks done for undefined variables
set -u

# INSTALL_TYPE=${1:-1}
# APP_INSTALL=${2:-1}
# SYSTEMD_INSTALL=${3:-1}

# if [ ${SYSTEMD_INSTALL} == "" ]; then
#   SYSTEMD_INSTALL="true"
# fi

# source "textutils.sh"

if [ "${HTTP_PROXY+x}" != "" ]; then
  export DOCKER_BUILD_ARGS="--build-arg http_proxy='${http_proxy}' --build-arg https_proxy='${https_proxy}' --build-arg HTTP_PROXY='${HTTP_PROXY}' --build-arg HTTPS_PROXY='${HTTPS_PROXY}' --build-arg NO_PROXY='localhost,127.0.0.1'"
  export DOCKER_RUN_ARGS="--env http_proxy='${http_proxy}' --env https_proxy='${https_proxy}' --env HTTP_PROXY='${HTTP_PROXY}' --env HTTPS_PROXY='${HTTPS_PROXY}' --env NO_PROXY='localhost,127.0.0.1'"
  export AWS_CLI_PROXY="export http_proxy='${http_proxy}'; export https_proxy='${https_proxy}'; export HTTP_PROXY='${HTTP_PROXY}'; export HTTPS_PROXY='${HTTPS_PROXY}'; export NO_PROXY='localhost,127.0.0.1';"
else
  export DOCKER_BUILD_ARGS=""
  export DOCKER_RUN_ARGS=""
  export AWS_CLI_PROXY=""
fi

# while (! docker stats --no-stream > /dev/null ); do
#   # Docker takes a few seconds to initialize
#   echo "Waiting for Docker to launch..." 2>&1 | tee -a ${CONSOLE_OUTPUT}
#   sleep 3
# done

if [ -d /opt/ace/bin ]; then
  msg="Skipping Installing App Docker Images and ACE"
  echo "$msg"
else
  msg="Installing App Docker Images..."
  echo "$msg"
  # docker run -d --privileged --name app-docker -v /var/lib/app-docker:/var/lib/docker -v /run:/opt/run docker:19.03.12-dind
  while (! docker exec -i $(docker ps | grep [_-]app-docker[_-] | awk '{print $1}') docker ps > /dev/null 2>&1 ); do 
    echo "Waiting for Docker to launch..." 
    sleep 0.5;
  done
  # docker exec -i app-docker sh -c 'docker -H unix:///opt/run/docker.sock save $(docker -H unix:///opt/run/docker.sock images --format "{{.Repository}}:{{.Tag}}" | grep glusterfs-plugin) | docker load'
  docker save $(docker images --format "{{.Repository}}:{{.Tag}}" | grep glusterfs-plugin) | docker exec -i $(docker ps | grep [_-]app-docker[_-] | awk '{print $1}') docker load
  # docker stop app-docker
  # docker rm app-docker
  DOCKER_ACE_CONSOLE=$(docker images | grep "[_-]console" | awk '{print $1}')
  docker run -it --rm --entrypoint="" -v /opt:/opt -v /var/lib/app-docker:/tmp/app-docker ${DOCKER_ACE_CONSOLE} rsync -a /tmp/app-docker/ /opt/ace/app-docker/
  docker pull docker:19.03.12
  docker pull docker:19.03.12-dind
  msg="Installing ACE Files..."
  echo "$msg"
  mkdir -p /opt/ace/
  rsync -rtc /ace/ /opt/ace/
fi

if [ ! -f /etc/ssl/dhparam.pem ]; then
  openssl dhparam -out /etc/ssl/dhparam.pem 2048
fi

if [ ${INSTALL_TYPE} == "demo" ]; then
  mkdir -p /etc/ssl/ace
  echo '[ "b8+87a00D33FD704a9deB1+DAb5B7Df917DFf7f2172=" ]' > /etc/ssl/ace/keyring.json
  msg="Installing certificates..."
  echo "$msg"
  cp /opt/ace/demo_certs/* /etc/ssl/
else
  cp -a /opt/ace/node_keys/* /etc/ssl/
fi

# if [ ${SYSTEMD_INSTALL} == "true" ]; then
#   msg="Installing Systemd service..."
#   echo "$msg"
#   cp /opt/acesystemd/ace.service /etc/systemd/system/
#   ln -s /etc/systemd/system/ace.service /etc/systemd/system/default.target.wants/ace.service
#   echo ""
#   echo ""
#   echo ""
#   msg="Run systemctl start ace"
#   echo "$msg"
#   echo ""
# fi
