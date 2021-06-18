#!/bin/bash

workdir=$(cd "$(dirname "$0")"; pwd)
source "${workdir}"/common.sh

check() {
  taskID=$1
  local ccDirNotExist=false
  local gitDirNotExist=false
  if [[ ! $(ls -A ${ccTmpRootPath}) ]] || [[ ! $(ls -A "${ccTmpRootPath}"/*_"${taskID}") ]]; then
    ccDirNotExist=true
  fi
  if [[ ! $(ls -A ${gitTmpRootPath}) ]] || [[ ! $(ls -A "${gitTmpRootPath}"/*_"${taskID}") ]]; then
    gitDirNotExist=true
  fi
  if ${ccDirNotExist} && ${gitDirNotExist} ; then
    exit 1
  fi
  exit 0
}

check "$1"
