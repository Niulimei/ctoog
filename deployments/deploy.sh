#!/bin/bash

######
#脚本名称：deploy.sh
#作用：完成server和worker的自动化安装部署
#传参说明：共需10个参数，依次分别为：
#ip地址,用户,密码,版本文件,工作目录,标志,server ip地址,server端口,worker ip地址,worker端口
######

work_path=$(
  cd "$(dirname "$0")" || exit 1
  pwd
)
scpParam="-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null"
scp="scp -r ${scpParam}"
ssh="ssh ${scpParam}"

copyFile() {
  local ip="$1"
  local user="$2"
  local password="$3"
  local versionFile="$4"
  local workDir="$5"
  rpm -ivh "${work_path}"/package/sshpass-*.rpm
  cd "${work_path}" || exit 1
  sshpass -p "${password}" ${scp} "${versionFile}" "${user}"@"${ip}":"${workDir}"
}

install() {
  local ip="$1"
  local user="$2"
  local password="$3"
  local versionFile="$4"
  local workDir="$5"
  local flag="$6"
  local sip="$7"
  local sport="$8"
  local wip="$9"
  local wport="${10}"
  local cmd0 cmd1 cmd2 cmd3
  cmd0="ps -ef|grep translator|grep -v grep|awk '{print $2}'|xargs -n1 kill -9 &>/dev/null"
  cmd1="cd ${workDir};mkdir -p ${flag};rm -rf ./${flag}/*;tar zxf $(basename "${versionFile}") -C ${flag}"
  if [[ ${flag} == "translator-server" ]]; then
    cmd2="cd ${flag};rm -rf script *-worker*;sed -i 's#SERVER_PORT#${sport}#' ${flag}.yaml;sed -i 's#SERVER_IP#${sip}#' ${flag}.yaml"
  elif [[ ${flag} == "translator-worker" ]]; then
    cmd2="cd ${flag};rm -rf frontend sql *-server*;sed -i 's#SERVER_PORT#${sport}#' ${flag}.yaml;sed -i 's#SERVER_IP#${sip}#' ${flag}.yaml;sed -i 's#WORKER_PORT#${wport}#' ${flag}.yaml;sed -i 's#WORKER_IP#${wip}#' ${flag}.yaml"
  fi
  cmd3="./${flag} start;./${flag} status"
  sshpass -p "${password}" ${ssh} "${user}"@"${ip}" "${cmd0};${cmd1};${cmd2};${cmd3}"
}

main() {
  copyFile "$1" "$2" "$3" "$4" "$5"
  install "$1" "$2" "$3" "$4" "$5" "$6" "$7" "$8" "$9" "${10}"
}

main "$1" "$2" "$3" "$4" "$5" "$6" "$7" "$8" "$9" "${10}"
