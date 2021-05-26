#!/bin/bash

workdir=$(cd $(dirname $0); pwd)
source "${workdir}"/common.sh

clean() {
  taskID=$1
  gitDirs=$(ls ${gitTmpRootPath} | grep -e ".*_${taskID}$")
  for gitDir in ${gitDirs}; do
    rm -rf ${gitTmpRootPath}/${gitDir}
  done
}

main() {
  clean $1
}

main $1
