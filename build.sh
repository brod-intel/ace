#!/bin/bash

# we want to have some checks done for undefined variables
set -u

source "scripts/textutils.sh"

if [ -z "${TOKEN:-}" ]; then 
  echo "Please set global vavriable TOKEN with 'export TOKEN='"; 
  exit 1
fi

if [ "${HTTP_PROXY+x}" != "" ]; then
  export DOCKER_BUILD_ARGS="--build-arg http_proxy=${http_proxy} --build-arg https_proxy=${https_proxy} --build-arg HTTP_PROXY=${HTTP_PROXY} --build-arg HTTPS_PROXY=${HTTPS_PROXY} --build-arg NO_PROXY=${NO_PROXY} --build-arg no_proxy=${no_proxy}"
  export DOCKER_RUN_ARGS="--env http_proxy=${http_proxy} --env https_proxy=${https_proxy} --env HTTP_PROXY=${HTTP_PROXY} --env HTTPS_PROXY=${HTTPS_PROXY} --env NO_PROXY=${NO_PROXY} --env no_proxy=${no_proxy}"
  export AWS_CLI_PROXY="export http_proxy=${http_proxy}; export https_proxy=${https_proxy}; export HTTP_PROXY=${HTTP_PROXY}; export HTTPS_PROXY=${HTTPS_PROXY}; export NO_PROXY=${NO_PROXY}; export no_proxy=${no_proxy}"
else
  export DOCKER_BUILD_ARGS=""
  export DOCKER_RUN_ARGS=""
  export AWS_CLI_PROXY=""
fi

# Build GO serf handlers
msg="Building GO Serf handlers, this can take a few minutes..."
printBanner "$msg"
logMsg "$msg"
./scripts/buildSerfHandlers.sh

# Build Autonomous Cluster for the Edge
msg="Building ACE, this can take a few minutes..."
printBanner "$msg"
logMsg "$msg"
source "scripts/buildACE.sh"

# Build Dynamic Hardware/Software Orchestrator
msg="Building DXO, this can take a few minutes..."
printBanner "$msg"
logMsg "$msg"
source "scripts/buildDXO.sh"

msg="Building Complete"
printBanner "$msg"
