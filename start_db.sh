#!/bin/bash

# Stop the container if it's already running
docker stop rocket-db

# Remove the stopped container so we can start a new one
docker rm rocket-db

# Run the container. Environment variables are populated from .db.env and
# determine the database name, user, and password. The container is attached
# to a user-defined network so that it can be reached from the app, and bound
# to a volume so database contents persist across restarts.
docker run \
    --name rocket-db \
    --env-file /go/src/github.com/ubclaunchpad/rocket/.db.env \
    --network rocket-net \
    -p 5432:5432 \
    -v pgdata:/var/lib/postgresql/data \
    -d \
    postgres:9.6
