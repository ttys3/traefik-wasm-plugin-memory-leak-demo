#!/usr/bin/env bash

set -eou pipefail

pushd ./tool/
go build .
popd

./tool/update-dyn-config -sleep 500ms -count 100

