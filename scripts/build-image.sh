#!/bin/bash

# set -x

WDIR=$(dirname ${BASH_SOURCE[0]})
source $WDIR/settings.sh

DOCKERFILE="./docker/Dockerfile"

(
    cd $WDIR/..
    docker build $TAGS --file $DOCKERFILE .
)

