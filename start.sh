#!/bin/bash

# Wait for postgres instance to boot before starting Rocket.

set -e
echo "Waiting for Postgres to start..."

while true; do
  if eval "pg_isready -h postgres -p 5432"; then
    break
  fi
  sleep 5
done

echo "Postgres ready"
eval rocket