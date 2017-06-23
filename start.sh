#!/bin/sh
# Starts rocket by populating the environment and running
source $ROCKET_PATH/.env
echo $ROCKET_HOST
$GOPATH/bin/rocket