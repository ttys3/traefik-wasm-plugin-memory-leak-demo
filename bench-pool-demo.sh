#!/usr/bin/env bash

set -eou pipefail

ab -v 3 -k -n 10000 -c 1000 http://localhost:8079/

