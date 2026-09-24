#!/bin/sh
set -eu

umask 022

mkdir -p /app/data /var/lib/grok2api-quality-guard
chmod 755 /app/data /var/lib/grok2api-quality-guard

exec /app/grok2api --config /app/config.yaml --listen 0.0.0.0:8000
