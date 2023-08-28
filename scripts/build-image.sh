#!/bin/bash

# set -x

WDIR=$(dirname ${BASH_SOURCE[0]})

source $WDIR/settings.sh

DOCKERFILE="./docker/Dockerfile"
DOCKERFILE_ALPINE="./docker/Dockerfile-alpine"

# Format the files

(
    cd $WDIR/..
    echo Formatting go files
    go fmt -x ./...
)

(
    set_version $(get_version_from_go_run)
    echo Building version $VERSION
    cd $WDIR/..
    docker build $(get_all_tags) --file $DOCKERFILE .
)

