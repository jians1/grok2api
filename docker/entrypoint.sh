#!/bin/sh
set -eu

umask 077

quality_guard_dir=/var/lib/grok2api-quality-guard
mkdir -p "${quality_guard_dir}" /app/data
chown grok2api:grok2api "${quality_guard_dir}" /app/data
chmod 0700 "${quality_guard_dir}"

config_source="${GROK2API_CONFIG_SOURCE:-/run/grok2api/config.yaml}"
if [ -f "${config_source}" ]; then
  cp "${config_source}" /app/config.yaml
elif [ ! -f /app/config.yaml ]; then
  # 无挂载配置时自动生成最小启动配置（与 ensure.go 默认值对齐）。
  cat > /app/config.yaml <<'EOF'
frontend:
  publicApiBaseURL: "http://127.0.0.1:8000"
  staticPath: "/app/frontend/dist"
database:
  driver: sqlite
  sqlite:
    path: "./data/backend.db"
media:
  driver: local
  local:
    path: "./data/media"
EOF
fi

chown grok2api:grok2api /app/config.yaml
chmod 0600 /app/config.yaml

exec su-exec grok2api:grok2api "$@"
