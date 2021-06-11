#!/bin/bash

cp="/bin/cp -rf"

main() {
  mkdir -p version;rm -rf ./version/*
  cd ../cmd/translator-server || exit 1
  go build -ldflags "-w -s" -o ../../deployments/version/translator-server
  ${cp} translator-server.yaml ../../deployments/version
  cd - || exit 1
  cd ../cmd/translator-worker || exit 1
  go build -ldflags "-w -s" -o ../../deployments/version/translator-worker
  ${cp} translator-worker.yaml ../../deployments/version
  cd - || exit 1
  ${cp} ../frontend ./version
  ${cp} ../script ./version
  ${cp} ../sql ./version
  cd version || exit 1
  tar zcf version_"$(date '+%Y%m%d%H%M%S')".tar.gz * --remove-files
}

main
