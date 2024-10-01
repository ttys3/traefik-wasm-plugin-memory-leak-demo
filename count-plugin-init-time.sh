#!/usr/bin/env bash

set -eou pipefail

watch -c -n2 "grep -c 'Demo load config' /tmp/a.log"


