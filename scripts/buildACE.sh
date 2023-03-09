#!/bin/bash

set -u

cd dockerfiles/

run "(1/9) Building Docker Image ACE Alpine Console" \
	"docker build --rm ${DOCKER_BUILD_ARGS} -t edge/console-alpine:1.0 ./console-alpine/Dockerfile" \
	${LOG_FILE}

run "(2/9) Building Docker Image GlusterFS" \
	"docker build --rm ${DOCKER_BUILD_ARGS} -t edge/glusterfs:7 ./glusterfs-server/Dockerfile" \
	${LOG_FILE}
	
run "(3/9) Building Docker Image GlusterFS REST" \
	"docker build --rm ${DOCKER_BUILD_ARGS} -t edge/glusterfs-rest:7 ./glusterfs-rest/Dockerfile" \
	${LOG_FILE}
	
run "(4/9) Building Docker Image Serf" \
	"docker build --rm ${DOCKER_BUILD_ARGS} -t edge/serf:0.8.4 ./serf/Dockerfile" \
	${LOG_FILE}

run "(5/9) Building Docker Image Dynamic Hardware Orchestrator" \
	"docker build --rm ${DOCKER_BUILD_ARGS} -t edge/dho:1.0 ./dho/Dockerfile" \
	${LOG_FILE}

run "(6/9) Building Docker Image App Docker" \
	"docker build --rm ${DOCKER_BUILD_ARGS} -t edge/app-docker:1.0 ./appdocker/Dockerfile" \
	${LOG_FILE}

run "(7/9) Building Docker Image Random Number Generator" \
	"docker build --rm ${DOCKER_BUILD_ARGS} -t edge/rngd:1.0 ./rngd/Dockerfile.rngd" \
	${LOG_FILE}

run "(8/9) Building Docker Image Gluster Plugin" \
	"docker build --rm ${DOCKER_BUILD_ARGS} -t edge/glusterfs-plugin:1.0 ./glusterfs-plugin/Dockerfile" \
	${LOG_FILE}

run "(9/9) Building Docker Image ACE Core" \
	"docker build --rm ${DOCKER_BUILD_ARGS} -t edge/core:1.0 ./core/Dockerfile" \
	${LOG_FILE}

cd - > /dev/null
