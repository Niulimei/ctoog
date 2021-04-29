#!/bin/bash

set -e

initGitRepo(){
  repoUrl=$1
  branchName=$2
  tmpGitDir=$3

  masterBranchName="origin/init_master"
  mkdir -p ${tmpGitDir}
  cd ${tmpGitDir}
  git init
  git remote add origin ${repoUrl}
  git fetch --all
  remoteBrList=$(git branch -r)
  branchExist=false
  branchMasterExist=false
  for br in ${remoteBrList}; do
    if [[ $(basename ${br}) == ${branchName} ]]; then
      branchExist=true
    fi
    if [[ ${br} == ${masterBranchName} ]]; then
      branchMasterExist=true
    fi
  done
  if $branchExist; then
    git checkout -b ${branchName} origin/${branchName}
    git pull origin ${branchName}
    return
  fi
  if $branchMasterExist; then
    git pull origin init_master
  else
    git checkout -b init_master
    touch ./.init_master
    git add .init_master
    git commit -m "init master"
    rm -rf .init_master
    git add .init_master
    git commit -m "delete master init file"
    git push origin init_master
  fi
  git checkout -b ${branchName}
}

pullCCAndPush(){
  repoUrl=$1
  branchName=$2
  random=$(cat /proc/sys/kernel/random/uuid)
  local tmpGitDir="/home/tmp/git/${random}"
  local tmpCCDir="/home/tmp/pvobs_view/${random}"
  pvobList=$(cleartool lsvob -s | grep pvob)
  for pv in ${pvobList}; do
    streamList=$(cleartool lsstream -s -invob ${pv})
    for sm in ${streamList}; do
      compList=$(cleartool lsstream -fmt %[components]p ${sm}@${pv})
      for comp in ${compList}; do
        mkdir -p ${tmpCCDir}
        cleartool mkview -snapshot -tag ${pv}_${sm}_${comp} -stgloc -auto -stream ${sm}@${pv} ${tmpDir}
        cd ${tmpCCDir}
        cleartool update -add_loadrules ${comp}
        initGitRepo ${repoUrl} ${branchName} ${tmpGitDir}
        cp -rf ${tmpCCDir}/* ${tmpGitDir}/*
        git add .
        git commit -m "import from cc,first commit $(date '+%Y%m%d%H%M%S')"
        git push origin ${branchName}
#        cleartool rmview ${tmpCCDir}
      done
    done
  done
}
