#!/bin/bash

WDIR=$(dirname ${BASH_SOURCE[0]})
source $WDIR/settings.sh

get_remote_tags ; echo
get_local_tags ; echo ; echo

get_all_tags ; echo