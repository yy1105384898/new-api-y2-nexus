#!/bin/sh
set -e

if [ "$(id -u)" = "0" ]; then
    mkdir -p /data /app/logs
    chown -R newapi:newapi /data /app/logs 2>/dev/null || true
    exec su-exec newapi "$0" "$@"
fi

if [ "${1#-}" != "$1" ]; then
    set -- /new-api "$@"
fi

exec "$@"
