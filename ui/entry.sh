#!/bin/sh
cat <<EOT > /usr/share/nginx/html/config.js
{
	"domain": "$AUTH0_DOMAIN",
	"client_id": "$AUTH0_CLIENT_ID",
	"audience": "vueapp"
}
EOT
