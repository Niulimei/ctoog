#!/bin/bash

git_repo_url=$1

# 检查git仓库是否存在以及git用户密码是否正确并有权限
err=$(git ls-remote "${git_repo_url}" 2>&1)
if [[ ! $? -eq 0 ]]; then
  echo "${err}"
  exit 1
fi

exit 0
