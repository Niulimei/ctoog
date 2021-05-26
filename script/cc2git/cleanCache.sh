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
    cleartool rmview ${ccTmpRootPath}/${ccDir}
  done
  gitDirs=$(ls ${gitTmpRootPath} | grep -e ".*_${taskID}$")
  for gitDir in ${gitDirs}; do
    if [[ ! ${exception} == "" ]]; then
      if [[ ${gitDir} =~ ${exception} ]]; then
        continue
      fi
    fi
    rm -rf ${gitTmpRootPath}/${gitDir}
  done
}

main() {
  clean $1 $2
}

main $1 $2
