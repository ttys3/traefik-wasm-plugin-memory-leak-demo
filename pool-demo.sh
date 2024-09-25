#!/usr/bin/env bash

set -eou pipefail

 env GOGC=50 GODEBUG=gctrace=1 go run ./go-sync-pool.go

