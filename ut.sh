#!/bin/bash

modules=(
  badger
  reaper
)

go clean -cache

for module in ${modules[*]}; do
  cd ${module}
  go clean -cache
  go test . -v
  cd -
done
