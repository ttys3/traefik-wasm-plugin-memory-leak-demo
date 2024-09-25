#!/usr/bin/env bash

set -eou pipefail

sudo watch -n1 'ps_mem  | grep -E "traefik|go-sync-pool"'


