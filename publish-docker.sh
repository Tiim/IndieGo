#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail


docker build -t tiimb/indiego:latest .

docker push tiimb/indiego:latest