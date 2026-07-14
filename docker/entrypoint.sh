#!/bin/sh
set -eu

umask 022

mkdir -p /app/data
chmod 755 /app/data

exec /app/grok2api --config /app/config.yaml --listen 0.0.0.0:8000