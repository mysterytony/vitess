#!/usr/bin/env bash

set -e

printf "\nStopping Vitess cluster\n"

export VTROOT=/vagrant
export VTDATAROOT=/tmp/vtdata-dev
export MYSQL_FLAVOR=MySQL56
cd "$VITESS_WORKSPACE"/examples/local

./vtgate-down.sh
UID_BASE='100' ./vttablet-down.sh
UID_BASE='200' ./vttablet-down.sh
UID_BASE='300' ./vttablet-down.sh
UID_BASE='400' ./vttablet-down.sh
UID_BASE='500' ./vttablet-down.sh
UID_BASE='600' ./vttablet-down.sh
./vtctld-down.sh
./zk-down.sh

rm -rf $VTDATAROOT

printf "\nVitess cluster stopped successfully.\n\n"
