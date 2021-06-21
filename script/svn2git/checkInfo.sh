#!/bin/bash

svn_username=$1
svn_password=$2
svn_url=$3
git_repo_url=$4

# 检查svn仓库是否存在以及svn用户密码是否正确并有权限
err=$(svn ls --non-interactive --username "${svn_username}" --password "${svn_password}" "${svn_url}" 2>&1)
if [[ ! $? -eq 0 ]]; then
  echo "${err}"
  exit 1
fi

# 检查git仓库是否存在以及git用户密码是否正确并有权限
err=$(git ls-remote "${git_repo_url}" 2>&1)
if [[ ! $? -eq 0 ]]; then
  echo "${err}"
  exit 1
fi

exit 0
