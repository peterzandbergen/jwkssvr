#!/bin/bash

WDIR=$(dirname ${BASH_SOURCE[0]})
source $WDIR/settings.sh

docker push $REMOTE_IMAGE