#!/bin/sh
set -e

KEY_PATH="/app/data/.ssh/id_ed25519"

if [ ! -f "$KEY_PATH" ]; then
  echo "🔑 Generating SSH host key..."
  mkdir -p "$(dirname "$KEY_PATH")"
  ssh-keygen -t ed25519 -f "$KEY_PATH" -N "" -C "ssh-portal-host-key"
  echo "✅ Host key generated at $KEY_PATH"
else
  echo "✅ Host key already exists."
fi

echo "🚀 Starting SSH Portal on :2222..."
exec /app/ssh-portal
