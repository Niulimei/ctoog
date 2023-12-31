#!/bin/bash

######
#脚本名称：cc2git.sh
#作用：完成CC代码的拉取，git仓库的初始化，代码向git的推送
#传参说明：共需10个参数，依次分别为：
#pvob名称，component名称，stream名称，gitRepoURL，git目标分支代码，任务ID，是否保留空目录(是：true，否：false)，用户名，邮箱, 空目录占位文件名称
######

export LANG="zh_CN.UTF-8"
set -e
workdir=$(cd $(dirname $0); pwd)
source ${workdir}/common.sh

initGitRepo(){
  echo "Initializing git repository..."
  repoUrl=$1
  branchName=$2
  tmpGitDir=$3
  username=$4
  email=$5
  masterBranchName="origin/init_master"
  mkdir -p ${tmpGitDir}
  cd ${tmpGitDir}
  git init
  git config --local core.longpaths true
  git config user.name "${username}"
  git config user.email "${email}"
  git config push.default simple
  git remote add origin ${repoUrl}
  git remote update
  git fetch --all
  git fetch -p origin
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
    git pull
    return
  fi
  if $branchMasterExist; then
    git checkout -b origin/init_master
  else
    git checkout -b init_master
    touch ./.init_master
    git add -A .init_master
    git commit --allow-empty -m "init master"
    rm -rf .init_master
    git add -A .init_master
    git commit --allow-empty -m "delete master init file"
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
  username=$8
  email=$9
  emptyFileName=${10}
  combainNameAdapt=$(echo -n ${pvobName}_${streamName}_${componentName} | sed 's/\//_/g')
  local tmpGitDir="${gitTmpRootPath}/${combainNameAdapt}_${taskID}"
  local tmpCCDir="${ccTmpRootPath}/${combainNameAdapt}_${taskID}"
  local tmpCCDirExist=false
  local tmpGitDirExist=false
  echo "Cloning code..."
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
    rm -rf ${tmpGitDir}
    tmpGitDirExist=true
  fi
  initGitRepo ${gitRepoUrl} ${gitBranchName} ${tmpGitDir} ${username} ${email}
  rm -rf ${tmpGitDir:?}/*
  cd ${tmpGitDir}
  echo "Copying files..."
  cp -rf ${tmpCCDir}${componentName}/* ${tmpGitDir}/
  if [[ ${containEmptyDir} == "true" ]]; then
    find ${tmpGitDir} -type d -empty -not -path "./.git/*" -exec touch {}/"${emptyFileName}" \;
  fi
  bash ${workdir}/changeCharSet.sh ${tmpGitDir}
  git add -A .
  echo "Pushing code..."
  if $tmpCCDirExist && $tmpGitDirExist; then
    lastMessage=$(git status | tail -n 2)
    #nothing to commit, working tree clean
    noCommit='nothing to commit'
    if [[ $lastMessage =~ $noCommit ]]; then
      set +e
      git push origin ${gitBranchName}
      set -e
    else
      git commit --allow-empty -m "sync from cc, update commit $(date '+%Y%m%d%H%M%S')"
      git push origin ${gitBranchName}
    fi
  else
    git commit --allow-empty -m "sync from cc, first commit $(date '+%Y%m%d%H%M%S')"
    git push origin ${gitBranchName}
  fi
}

main(){
  mkdir -p ${ccTmpRootPath}
  mkdir -p ${gitTmpRootPath}
  pullCCAndPush $1 $2 $3 $4 $5 $6 $7 $8 $9 ${10}
}

main $1 $2 $3 $4 $5 $6 $7 $8 $9 ${10}
