#!/bin/bash

WDIR=$(dirname ${BASH_SOURCE[0]})
source $WDIR/settings.sh

set_version $(get_version_from_go_run)

for i in "${REMOTE_TAGS[@]}"
do
    echo Pushing $i
    docker push $i
done