#!/usr/bin/env bash

TIMEOUT=5
HOST=$1
PORT=$2

echo "Waiting for $HOST:$PORT to be available..."

# Loop until nc (netcat) returns that the port is open
until nc -z -v -w $TIMEOUT "$HOST" "$PORT"; do
  >&2 echo "Waiting for $HOST:$PORT to be available..."
  sleep 1
done

echo "$HOST:$PORT is available!"
exec "${@:3}"
