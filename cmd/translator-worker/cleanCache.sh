#!/bin/bash

workdir=$(cd $(dirname $0); pwd)
source ${workdir}/common.sh

clean() {
  taskID=$1
  exception=$2
  ccDirs=$(ls ${ccTmpRootPath} | grep -e ".*_${taskID}$")
  for ccDir in ${ccDirs}; do
    if [[ ! ${exception} == "" ]]; then
      if [[ ${ccDir} =~ ${exception} ]]; then
        continue
      fi
    fi
    set -e
    cleartool rmview ${ccTmpRootPath}/${ccDir}
    set +e
  done
  gitDirs=$(ls ${gitTmpRootPath} | grep -e ".*_${taskID}$")
  for gitDir in ${gitDirs}; do
    if [[ ! ${exception} == "" ]]; then
      if [[ ${gitDir} =~ ${exception} ]]; then
        continue
      fi
    fi
    set -e
    rm -rf ${gitTmpRootPath}/${gitDir}
    set +e
  done
}

main() {
  clean $1 $2
}

main $1 $2
