#!/bin/bash

######
#脚本名称：cc2git_onlygit.sh
#作用：完成git仓库的初始化，代码向git的推送
#传参说明：共需11个参数，依次分别为：
#gitRepoURL，git目标分支代码，任务ID，是否保留空目录(是：true，否：false)，用户名，邮箱, 空目录占位文件名称, gitignore文件内容, cc source dir, component名称, git tmp dir
######

set -x
export LANG="zh_CN.UTF-8"
set -e

initGitRepo(){
  echo "Initializing git repository..."
  repoUrl=$1
  branchName=$2
  tmpGitDir=$3
  username=$4
  email=$5
  masterBranchName="origin/init_master"
  mkdir -p "${tmpGitDir}" -m 777
  cd "${tmpGitDir}"
  git init
  git config --local core.longpaths true
  git config user.name "${username}"
  git config user.email "${email}"
  git config push.default simple
  git remote add origin "${repoUrl}"
  git remote update
  git fetch --all
  git fetch -p origin
  remoteBrList=$(git branch -r)
  branchExist=false
  branchMasterExist=false
  for br in ${remoteBrList}; do
    if [[ "$(basename "${br}")" == "${branchName}" ]]; then
      branchExist=true
    fi
    if [[ "${br}" == "${masterBranchName}" ]]; then
      branchMasterExist=true
    fi
  done
  if $branchExist; then
    git checkout -b "${branchName}" origin/"${branchName}"
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
  git checkout -b "${branchName}"
}

pullCCAndPush(){
  gitRepoUrl=$1
  gitBranchName=$2
  taskID=$3
  containEmptyDir=$4
  username=$5
  email=$6
  emptyFileName=$7
  gitignoreContent=$8
  ccDirName="$9"
  combineNameAdapt=$(echo -n "${ccDirName}" | sed 's/\//_/g')
  componentName="${10}"
  local tmpGitDir="${gitTmpRootPath}/${combineNameAdapt}_${taskID}"
  local tmpCCDir="${ccDirName}"
  local tmpCCDirExist=false
  local tmpGitDirExist=false

  if [[ -d "${tmpGitDir}" ]]; then
    rm -rf "${tmpGitDir}"
    tmpGitDirExist=true
  fi
  initGitRepo "${gitRepoUrl}" "${gitBranchName}" "${tmpGitDir}" "${username}" "${email}"
  tmpCCDirExist=true
  rm -rf "${tmpGitDir:?}"/*
  cd "${tmpGitDir}"
  echo "Copying files..."
  cp -rf "${tmpCCDir}"/"${componentName}"/. "${tmpGitDir}"/ >/dev/null
  if [[ ${containEmptyDir} == "true" ]]; then
    find "${tmpGitDir}" -type d -empty -not -path "./.git/*" -exec touch {}/"${emptyFileName}" \;
  fi
#  bash "${workdir}"/changeCharSet.sh "${tmpGitDir}" &>/dev/null
  if [[ -n "${gitignoreContent}" ]]; then
    echo -e "${gitignoreContent}" >./.gitignore
  else
    rm -rf ./.gitignore
  fi
  git add -A .
  echo "Pushing code..."
  if $tmpCCDirExist && $tmpGitDirExist; then
    lastMessage=$(git status | tail -n 2)
    noCommit='nothing to commit'
    if [[ $lastMessage =~ $noCommit ]]; then
      set +e
      git push origin "${gitBranchName}"
      set -e
    else
      git commit --allow-empty -m "sync from cc, update commit $(date '+%Y%m%d%H%M%S')" >/dev/null
      git push origin "${gitBranchName}"
    fi
  else
    git commit --allow-empty -m "sync from cc, first commit $(date '+%Y%m%d%H%M%S')" >/dev/null
    git push origin "${gitBranchName}"
  fi
}

main(){
  if [[ -n "${11}" ]]; then
    gitTmpRootPath="${11}"
  else
    gitTmpRootPath="/home/tmp/git"
  fi
  mkdir -p "${gitTmpRootPath}" -m 777
  pullCCAndPush "$1" "$2" "$3" "$4" "$5" "$6" "$7" "$8" "$9" "${10}" "${11}"
}

main "$1" "$2" "$3" "$4" "$5" "$6" "$7" "$8" "$9" "${10}" "${11}"
