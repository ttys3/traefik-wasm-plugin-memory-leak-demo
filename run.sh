#!/usr/bin/env bash

set -eou pipefail

truncate -s 0 /tmp/a.log
env GOGC=50 GODEBUG=gctrace=1 ./traefik --configFile=./traefik.yaml 2>&1 | tee /tmp/a.log

