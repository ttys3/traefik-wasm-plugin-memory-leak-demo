#!/usr/bin/env bash

set -eou pipefail

for i in {1..10000}
do
   curl -vvv http://localhost:6688/
   sleep 0.1
done

