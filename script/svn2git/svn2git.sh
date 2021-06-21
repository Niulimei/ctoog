#!/bin/bash

######
#脚本名称：svn2git.sh
#作用：完成SVN代码的拉取，git仓库的初始化，代码向git的推送
#传参说明：共需8个参数，依次分别为：
#svnRepoURL，gitRepoURL，任务ID，是否保留空目录(是：true，否：false)，用户名，邮箱, 空目录占位文件名称, 用户名称映射文件, svn用户名, svn密码
######

set -x
export LANG="zh_CN.UTF-8"
set -e
workdir=$(cd "$(dirname "$0")"; pwd)
source "${workdir}"/common.sh

configGitRepo(){
  echo "Initializing git repository..."
  repoUrl=$1
  tmpGitDir=$2
  username=$3
  email=$4
  cd "${tmpGitDir}"
  git config --local core.longpaths true
  git config user.name "${username}"
  git config user.email "${email}"
  git config push.default simple
  git remote add origin "${repoUrl}"
  git remote update
  git fetch --all
  git fetch -p origin
}

pullCCAndPush(){
  svnRepoUrl=$1
  gitRepoUrl=$2
  taskID=$3
  containEmptyDir=$4
  username=$5
  email=$6
  if [[ -z ${email} ]];then
    email="default@empty.com"
  fi
  emptyFileName=$7
  userFile=$8
  gitignoreContent=$9
  svnUser="${10}"
  svnPassword="${11}"
  combineNameAdapt=$(basename "${svnRepoUrl}")
  local tmpGitDir="${gitTmpRootPath}/${combineNameAdapt}_${taskID}"
  local tmpGitDirExist=false
  echo "Cloning code..."
  if [[ -d ${tmpGitDir} ]]; then
    rm -rf "${tmpGitDir}"
    tmpGitDirExist=true
  fi
  if [[ -f ${userFile} ]]; then
    echo "${svnPassword}" | git svn clone --username "${svnUser}" --authors-file="${userFile}" --no-metadata --prefix "" "${svnRepoUrl}" "${tmpGitDir}" >/dev/null
  else
    echo "${svnPassword}" | git svn clone --username "${svnUser}" --no-metadata --prefix "" "${svnRepoUrl}" "${tmpGitDir}" >/dev/null
  fi
  rm -rf "${userFile}"
  configGitRepo "${gitRepoUrl}" "${tmpGitDir}" "${username}" "${email}"
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
  if $tmpGitDirExist; then
    lastMessage=$(git status | tail -n 2)
    #nothing to commit, working tree clean
    noCommit='nothing to commit'
    if [[ $lastMessage =~ $noCommit ]]; then
      set +e
      git push origin --mirror
      set -e
    else
      git commit --allow-empty -m "sync from svn, update commit $(date '+%Y%m%d%H%M%S')" >/dev/null
      git push origin --mirror
    fi
  else
    git commit --allow-empty -m "sync from svn, first commit $(date '+%Y%m%d%H%M%S')" >/dev/null
    git push origin --mirror
  fi
}

main(){
  mkdir -p "${gitTmpRootPath}" -m 777
  pullCCAndPush "$1" "$2" "$3" "$4" "$5" "$6" "$7" "$8" "$9" "${10}" "${11}"
}

main "$1" "$2" "$3" "$4" "$5" "$6" "$7" "$8" "$9" "${10}" "${11}"
