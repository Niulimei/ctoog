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
  git config --global http.postBuffer 2048M
  git config --global http.maxRequestBuffer 1024M
  git config --global core.compression 9
  git config --global ssh.postBuffer 2048M
  git config --global ssh.maxRequestBuffer 1024M
  git config --global pack.windowMemory 256m
  git config --global pack.packSizeLimit 256m
  git config --local core.longpaths true
  git config user.name "${username}"
  git config user.email "${email}"
  git config push.default simple
  git remote remove origin
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
  branchInfo="${12}"
  combineNameAdapt=$(basename "${svnRepoUrl}")
  local tmpGitDir="${gitTmpRootPath}/${combineNameAdapt}_${taskID}"
  local tmpGitDirExist=false
  local tmpGitProj=`echo $PROJECT_KEY | awk '{print tolower($0)}'`
  local tmpGitSlug=`echo ${combineNameAdapt}_${taskID} | awk '{print tolower($0)}'`
  echo $tmpGitProj
  echo $tmpGitSlug
  curl -v -u "$BITBUCKET_USERNAME:$BITBUCKET_PASSWORD" --request POST \
  --url 'http://'$BITBUCKET_HOST'/rest/api/latest/projects/'$PROJECT_KEY'/repos' \
  --header 'Accept: application/json' \
  --header 'Content-Type: application/json' \
  --data '{
  "name": "'$tmpGitSlug'",
  "slug": "'$tmpGitSlug'"
}'
  echo "Cloning code..."
  if [[ -d ${tmpGitDir} ]]; then
    rm -rf "${tmpGitDir}"
    tmpGitDirExist=true
  fi
  rm -rf /root/.subversion/auth
  userFileInfo=`cat "${userFile}"`
  CONFIGURE=$(cat <<END
   {
     "url" : "${svnRepoUrl}",
     "credentials" : {
         "username" : "${svnUser}",
         "password" : "${svnPassword}"
     },
     "layout" : {
         "type" : "MANUAL",
         "branches" : ["${branchInfo}"]
     },
     "config" : {
         "svn.fetchInterval" : 0
     },
     "authors" : ${userFileInfo}
   }
END
)
  echo "$CONFIGURE"
  curl -v -u "$BITBUCKET_USERNAME:$BITBUCKET_PASSWORD" \
  -H "Content-Type: application/json" \
  -H "X-Atlassian-Token:no-check" \
  -X POST \
  --data "$CONFIGURE" \
  'http://'$BITBUCKET_HOST'/rest/svn/1.0/projects/'$PROJECT_KEY'/repos/'$tmpGitSlug'/configure?start=import&async=false'

  git clone 'http://'$BITBUCKET_GIT_USER':'$BITBUCKET_GIT_PASSWORD'@'$BITBUCKET_GIT_HOST'/scm/'$tmpGitProj'/'$tmpGitSlug'.git' "$tmpGitDir"
  cd "${tmpGitDir}"
  echo "Pushing code..."
  configGitRepo "${gitRepoUrl}" "${tmpGitDir}" "${username}" "${email}"
  if $tmpGitDirExist; then
    lastMessage=$(git status | tail -n 2)
    #nothing to commit, working tree clean
    noCommit='nothing to commit'
    if [[ $lastMessage =~ $noCommit ]]; then
      set +e
      git push origin --all
      git push --tags
      set -e
    else
      git push origin --all
      git push origin --tags
    fi
  else
    git push origin --all
    git push --tags
  fi
  curl --request DELETE -v -u ""${BITBUCKET_USERNAME}":"${BITBUCKET_PASSWORD}"" \
  --url 'http://'$BITBUCKET_HOST'/rest/api/latest/projects/'$PROJECT_KEY'/repos/'$tmpGitSlug'' \
  --header 'Accept: application/json'
}

main(){
  mkdir -p "${gitTmpRootPath}" -m 777
  pullCCAndPush "$1" "$2" "$3" "$4" "$5" "$6" "$7" "$8" "$9" "${10}" "${11}" "${12}"
}

main "$1" "$2" "$3" "$4" "$5" "$6" "$7" "$8" "$9" "${10}" "${11}" "${12}"