#!/bin/bash

set -u

cd dockerfiles/

run "(1/9) Building Docker Image ACE Alpine Console" \
	"docker build --rm ${DOCKER_BUILD_ARGS} -t ace/console-alpine:1.0 -f ./console-alpine/Dockerfile ./console-alpine/" \
	../${LOG_FILE}

run "(2/9) Building Docker Image GlusterFS" \
	"docker build --rm ${DOCKER_BUILD_ARGS} -t ace/glusterfs:7 -f ./glusterfs-server/Dockerfile ./glusterfs-server/" \
	../${LOG_FILE}
	
run "(3/9) Building Docker Image GlusterFS REST" \
	"docker build --rm ${DOCKER_BUILD_ARGS} -t ace/glusterfs-rest:7 -f ./glusterfs-rest/Dockerfile ./glusterfs-rest/" \
	../${LOG_FILE}
	
run "(4/9) Building Docker Image Serf" \
	"docker build --rm ${DOCKER_BUILD_ARGS} -t ace/serf:0.8.4 -f ./serf/Dockerfile ./serf/" \
	../${LOG_FILE}

run "(5/9) Building Docker Image Dynamic Hardware Orchestrator" \
	"docker build --rm ${DOCKER_BUILD_ARGS} -t ace/dho:1.0 -f ./dho/Dockerfile ./dho/" \
	../${LOG_FILE}

run "(6/9) Building Docker Image App Docker" \
	"docker build --rm ${DOCKER_BUILD_ARGS} -t ace/app-docker:1.0 -f ./appdocker/Dockerfile ./appdocker/" \
	../${LOG_FILE}

run "(7/9) Building Docker Image Random Number Generator" \
	"docker build --rm ${DOCKER_BUILD_ARGS} -t ace/rngd:1.0 -f ./rngd/Dockerfile ./rngd/" \
	../${LOG_FILE}

run "(8/9) Building Docker Image Gluster Plugin" \
	"cp serf/gluster/updateconf ./glusterfs-plugin/ && docker build --rm ${DOCKER_BUILD_ARGS} -t ace/glusterfs-plugin:1.0 -f ./glusterfs-plugin/Dockerfile ./glusterfs-plugin/" \
	../${LOG_FILE}

run "(9/9) Building Docker Image ACE Core" \
	"docker build --rm ${DOCKER_BUILD_ARGS} -t ace/core:1.0 -f ./core/Dockerfile ./core/" \
	../${LOG_FILE}

cd - > /dev/null
