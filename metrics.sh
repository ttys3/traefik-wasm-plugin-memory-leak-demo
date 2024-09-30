#!/usr/bin/env bash

set -eou pipefail

watch -c -n2 "curl -Ssf http://127.0.0.1:8080/metrics | grep -E '^go_goroutines|^traefik_open_connections'"


