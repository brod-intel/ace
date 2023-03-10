#!/bin/bash

# we want to have some checks done for undefined variables
set -u

TYPE=${1:-1}
APP_INSTALL=${2:-1}
SYSTEMD_INSTALL=${3:-1}

if [ ${SYSTEMD_INSTALL} == "" ]; then
  SYSTEMD_INSTALL="true"
fi

source "scripts/textutils.sh"

if [ "${HTTP_PROXY+x}" != "" ]; then
  export DOCKER_BUILD_ARGS="--build-arg http_proxy='${http_proxy}' --build-arg https_proxy='${https_proxy}' --build-arg HTTP_PROXY='${HTTP_PROXY}' --build-arg HTTPS_PROXY='${HTTPS_PROXY}' --build-arg NO_PROXY='localhost,127.0.0.1'"
  export DOCKER_RUN_ARGS="--env http_proxy='${http_proxy}' --env https_proxy='${https_proxy}' --env HTTP_PROXY='${HTTP_PROXY}' --env HTTPS_PROXY='${HTTPS_PROXY}' --env NO_PROXY='localhost,127.0.0.1'"
  export AWS_CLI_PROXY="export http_proxy='${http_proxy}'; export https_proxy='${https_proxy}'; export HTTP_PROXY='${HTTP_PROXY}'; export HTTPS_PROXY='${HTTPS_PROXY}'; export NO_PROXY='localhost,127.0.0.1';"
else
  export DOCKER_BUILD_ARGS=""
  export DOCKER_RUN_ARGS=""
  export AWS_CLI_PROXY=""
fi

if [ ${TYPE} == "demo" ]; then
  mkdir -p /etc/ssl/ace
  echo '[ "b8+87a00D33FD704a9deB1+DAb5B7Df917DFf7f2172=" ]' > /etc/ssl/ace/keyring.json
  run "Installing certificates..." "cp demo_certs/* /etc/ssl/"   ${LOG_FILE}
else
  cp -a node_keys/* /etc/ssl/
fi

if [ ${APP_INSTALL} == "skip" ]; then
  msg="Skipping Installing App Docker Images"
  printBanner "$msg"
else
  msg="Installing App Docker Images..."
  printBanner "$msg"
  logMsg "$msg"
  docker run -d --privileged --name app-docker -v /var/lib/app-docker:/var/lib/docker -v /var/run:/opt/run docker:19.03.12-dind
  while (! docker exec -i app-docker docker ps > /dev/null 2>&1 ); do sleep 0.5; done
  docker exec -i app-docker sh -c 'docker -H unix:///opt/run/docker.sock save $(docker -H unix:///opt/run/docker.sock images --format "{{.Repository}}:{{.Tag}}" | grep glusterfs-plugin) | docker load'
  docker stop app-docker
  docker rm app-docker
  docker run -it --rm --entrypoint="" -v /opt:/opt -v /var/lib/app-docker:/tmp/app-docker edge/console-alpine:1.0 rsync -a /tmp/app-docker/ /opt/ace/app-docker/
  docker pull docker:19.03.12
  docker pull docker:19.03.12-dind
fi

if [ ${SYSTEMD_INSTALL} == "true" ]; then
  msg="Installing Systemd service..."
  printBanner "$msg"
  logMsg "$msg"
  cp systemd/ace.service /etc/systemd/system/
  ln -s /etc/systemd/system/ace.service /etc/systemd/system/default.target.wants/ace.service
  echo ""
  echo ""
  echo ""
  msg="Run systemctl start ace"
  printBanner "$msg"
  echo ""
fi
