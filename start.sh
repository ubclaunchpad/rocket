#!/bin/bash
# Starts rocket by populating the environment and running

source /go/src/github.com/ubclaunchpad/rocket/.env
/go/bin/rocket >> $ROCKET_LOGFILE
