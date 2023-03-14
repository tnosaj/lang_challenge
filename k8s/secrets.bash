#!/bin/bash
secret=$(echo -n "$1" | base64 -w0)
cat << EOF
apiVersion: v1
kind: Secret
metadata:
  name: lang-secrets
  namespace: lang
data:
  redis: $secret
EOF
