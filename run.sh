#!/usr/bin/env bash

set -eou pipefail

env GOGC=50 GODEBUG=gctrace=1 ./traefik --configFile=./traefik.yaml

