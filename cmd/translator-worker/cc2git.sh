#!/bin/bash

######
#脚本名称：cc2git.sh
#作用：完成CC代码的拉取，git仓库的初始化，代码向git的推送
#传参说明：共需7个参数，依次分别为：
#pvob名称，component名称，stream名称，gitRepoURL，git目标分支代码，任务ID，是否保留空目录(是：true，否：false)
######

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
  containEmptyDir=$7
  combainNameAdapt=$(echo -n ${pvobName}_${streamName}_${componentName} | sed 's/\//_/g')
  local tmpGitDir="${gitTmpRootPath}/${combainNameAdapt}_${taskID}"
  local tmpCCDir="${ccTmpRootPath}/${combainNameAdapt}_${taskID}"
  local tmpCCDirExist=false
  local tmpGitDir=false
  if [[ -d ${tmpCCDir} ]]; then
    tmpCCDirExist=true
    cd ${tmpCCDir}
    cleartool update .
  else
    cleartool mkview -snapshot -tag ${combainNameAdapt}_${taskID} -stgloc -auto -stream ${streamName}@${pvobName} ${tmpCCDir}
    cd ${tmpCCDir}
    cleartool update -add_loadrules .${componentName}
  fi
  if [[ -d ${tmpGitDir} ]]; then
    tmpGitDirExist=true
  else
    initGitRepo ${gitRepoUrl} ${gitBranchName} ${tmpGitDir}
  fi
  cp -rf ${tmpCCDir}/* ${tmpGitDir}/*
  if [[ ${containEmptyDir} == "true" ]]; then
    find ${tmpGitDir} -type d -empty -not -path "./.git/*" -exec touch {}/.gitkeep \;
  fi
  git add .
  if $tmpCCDirExist && $tmpGitDirExist; then
    git commit -m "import from cc,update commit $(date '+%Y%m%d%H%M%S')"
  else
    git commit -m "import from cc,first commit $(date '+%Y%m%d%H%M%S')"
  fi
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
  mkdir -p ${ccTmpRootPath}
  mkdir -p ${gitTmpRootPath}
  pullCCAndPush $1 $2 $3 $4 $5 $6 $7
  #postClean
}

main $1 $2 $3 $4 $5 $6 $7
