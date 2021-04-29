#!/bin/bash

set -e

ccTmpRootPath="/home/tmp/pvobs_view"
gitTmpRootPath="/home/tmp/git"

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
    if [[ $(basename ${br}) == "${branchName}" ]]; then
      branchExist=true
    fi
    if [[ ${br} == "${masterBranchName}" ]]; then
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
  pvobName=$1
  componentName=$2
  streamName=$3
  gitRepoUrl=$4
  gitBranchName=$5
  taskID=$6
  local tmpGitDir="${gitTmpRootPath}/${taskID}"
  local tmpCCDir="${ccTmpRootPath}/${taskID}"
  local tmpCCDirExist=false
  local tmpGitDir=false
  if [[ -d ${tmpCCDir} ]]; then
    tmpCCDirExist=true
    cd ${tmpCCDir}
    cleartool update .
  else
    mkdir -p ${tmpCCDir}
    cleartool mkview -snapshot -tag ${pvobName}_${streamName}_${componentName}_${taskID} -stgloc -auto -stream ${streamName}@${pvobName} ${tmpCCDir}
    cd ${tmpCCDir}
    cleartool update -add_loadrules ${componentName}
  fi
  if [[ -d ${tmpGitDir} ]]; then
    tmpGitDirExist=true
  else
    initGitRepo ${gitRepoUrl} ${gitBranchName} ${tmpGitDir}
  fi
  if $tmpCCDirExist && $tmpGitDirExist; then
    return
  fi
  cp -rf ${tmpCCDir}/* ${tmpGitDir}/*
  git add .
  git commit -m "import from cc,first commit $(date '+%Y%m%d%H%M%S')"
  git push origin ${branchName}
}

postClean(){
  rm -rf ${gitTmpRootPath:?}/*
  views=$(ls ${ccTmpRootPath})
  for view in ${views}; do
    cleartool rmview ${ccTmpRootPath}/${view}
  done
  rm -rf ${ccTmpRootPath:?}/*
}

main(){
  pullCCAndPush $1 $2 $3 $4 $5 $6
  #postClean
}

main $1 $2 $3 $4 $5 $6
