#!/bin/bash

# Build the container from the project directory and tag as latest.
docker build \
    -t rocket-app:latest \
    /go/src/github.com/ubclaunchpad/rocket

# Stop the old container, if it's running.
docker stop rocket-app

# Remove the old container, so we can start a new one.
docker rm rocket-app

# Run Rocket. Environment variables are populated from .app.env and define the
# app configuration. The container is attached to a user-defined network so it
# can reach the database (which is on the same network). Port 80 is bound to
# the host so we can reach the API.
docker run \
    --name rocket-app \
    --env-file /go/src/github.com/ubclaunchpad/rocket/.app.env \
    --network rocket-net \
    -p 80:80 \
    -d \
    rocket-app
