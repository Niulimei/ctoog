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
  taskID=$3
  local tmpGitDir="${gitTmpRootPath}/${taskID}"
  local tmpCCDir="${ccTmpRootPath}/${taskID}"
  pvobList=$(cleartool lsvob -s | grep pvob)
  for pv in ${pvobList}; do
    streamList=$(cleartool lsstream -s -invob ${pv})
    for sm in ${streamList}; do
      compList=$(cleartool lsstream -fmt %[components]p ${sm}@${pv})
      for comp in ${compList}; do
        local tmpCCDirExist=false
        local tmpGitDir=false
        if [[ -d ${tmpCCDir} ]]; then
          tmpCCDirExist=true
          cd ${tmpCCDir}
          cleartool update .
        else
          mkdir -p ${tmpCCDir}
          cleartool mkview -snapshot -tag ${pv}_${sm}_${comp}_${taskID} -stgloc -auto -stream ${sm}@${pv} ${tmpCCDir}
          cd ${tmpCCDir}
          cleartool update -add_loadrules ${comp}
        fi
        if [[ -d ${tmpGitDir} ]]; then
          tmpGitDirExist=true
        else
          initGitRepo ${repoUrl} ${branchName} ${tmpGitDir}
        fi
        if $tmpCCDirExist && $tmpGitDirExist; then
          continue
        fi
        cp -rf ${tmpCCDir}/* ${tmpGitDir}/*
        git add .
        git commit -m "import from cc,first commit $(date '+%Y%m%d%H%M%S')"
        git push origin ${branchName}
      done
    done
  done
}

main(){
  pullCCAndPush $1 $2 $3
#  rm -rf ${gitTmpRootPath}/*
#  views=$(ls ${ccTmpRootPath})
#  for view in ${views}; do
#    cleartool rmview ${ccTmpRootPath}/${view}
#  done
#  rm -rf ${ccTmpRootPath}/*
}

main $1 $2 $3
