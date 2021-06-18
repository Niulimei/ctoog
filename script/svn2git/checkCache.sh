#!/bin/bash

workdir=$(cd "$(dirname "$0")"; pwd)
source "${workdir}"/common.sh

check() {
  taskID=$1
  local gitDirNotExist=false
  if [[ $(ls -A ${gitTmpRootPath}) ]] || [[ $(ls "${gitTmpRootPath}"/*_"${taskID}") == "" ]]; then
    gitDirNotExist=true
  fi
  if ${gitDirNotExist} ; then
    exit 1
  fi
  exit 0
}

check "$1"
