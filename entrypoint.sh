#!/bin/sh
set -e

KEY_PATH="/app/data/.ssh/id_ed25519"

echo "📁 Data dir contents:"
ls -la /app/data/ 2>/dev/null || echo "(empty or missing)"
ls -la /app/data/.ssh/ 2>/dev/null || echo "(.ssh missing)"

if [ ! -f "$KEY_PATH" ]; then
  echo "🔑 Generating SSH host key..."
  mkdir -p "$(dirname "$KEY_PATH")"
  ssh-keygen -t ed25519 -f "$KEY_PATH" -N "" -C "ssh-portal-host-key"
  echo "✅ Host key generated"
else
  echo "✅ Host key exists, reusing it"
fi

echo "🚀 Starting SSH Portal on :2222..."
exec /app/ssh-portal