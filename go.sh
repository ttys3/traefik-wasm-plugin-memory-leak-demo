#!/usr/bin/env bash

for (( ; ; ))
do
    curl -X GET -vvv --connect-timeout 1 --max-time 2  http://localhost:6688/
done


