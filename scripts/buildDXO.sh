#!/bin/bash

if [ -z "${TOKEN}" ]; then 
  echo "Please set global vavriable TOKEN with 'export TOKEN='"; 
  exit 1
fi

docker run -i --rm --privileged ${DOCKER_RUN_ARGS} -e TOKEN="${TOKEN}" -v /var/run:/var/run docker:stable sh -c "\
  apk update && apk add --no-cache \ 
    bash \
    coreutils \
    git \
    ncurses \
    pv \
    tar \
    wget

  mkdir /work && \
  cd work && \
  git clone --depth=1 https://${TOKEN}:${TOKEN}@github.com/brod-intel/dxo.git dxo && \
  cd dxo && \
  docker build ${DOCKER_BUILD_ARGS} -t ace/dxo:1.0 -f build/Dockerfile ."
