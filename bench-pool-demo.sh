#!/usr/bin/env bash

set -eou pipefail

ab -v 3 -k -n 40000 -c 20000 http://localhost:8079/

