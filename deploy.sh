#!/usr/bin/env bash

set -e

REMOTE=comments.tiim.ch
export DOCKER_HOST="ssh://tim@$REMOTE"

docker-compose up --build -d