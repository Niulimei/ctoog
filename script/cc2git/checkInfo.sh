#!/bin/bash

cc_user=$1
cc_password=$2
git_repo_url=$3

# 检查cc用户是否存在
id -u "${cc_user}" &>/dev/null
if [[ ! $? -eq 0 ]]; then
  echo user "${cc_user}" does not exist.
  exit 1
fi

# 检查cc密码是否正确
echo "${cc_password}" | su - "${cc_user}" -c "echo 1" &>/dev/null
if [[ ! $? -eq 0 ]]; then
  echo the password of user "${cc_user}" is wrong.
  exit 1
fi

# 检查git仓库是否存在以及git用户密码是否正确并有权限
err=$(git ls-remote "${git_repo_url}" 2>&1)
if [[ ! $? -eq 0 ]]; then
  echo "${err}"
  exit 1
fi

exit 0
