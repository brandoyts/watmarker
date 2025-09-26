#!/bin/sh
set -e

echo "🔧 Replacing environment variables in Nginx config..."

envsubst '${API_BASE_URL}' \
  < /etc/nginx/conf.d/default.conf.template \
  > /etc/nginx/conf.d/default.conf

echo "✅ Applied API_BASE_URL=${API_BASE_URL} to Nginx config"
