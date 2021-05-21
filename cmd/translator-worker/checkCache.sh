#!/bin/bash

workdir=$(cd $(dirname $0); pwd)
source ${workdir}/common.sh

check() {
  local ccDirNotExist=false
  local gitDirNotExist=false
  if [[ ! -d ${ccTmpRootPath} ]] || [[ $(ls ${ccTmpRootPath}) == "" ]]; then
    ccDirNotExist=true
  fi
  if [[ ! -d ${gitTmpRootPath} ]] || [[ $(ls ${gitTmpRootPath}) == "" ]]; then
    gitDirNotExist=true
  fi
  if ${ccDirNotExist} && ${gitDirNotExist} ; then
    exit 1
  fi
  exit 0
}

check
