#!/usr/bin/env bash

if [ -f .env ]; then
    set -a
    source .env
    set +a
fi

goose postgres "$CONNECTION_STRING" down
