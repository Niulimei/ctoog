#!/bin/bash

clean() {
  databaseFile=$1
  if [[ ! -f ${databaseFile} ]]; then
    echo "Database File Error."
    exit 1
  fi
  for db in log task match_info task_log task_command_out worker schedule plan; do
    echo "Cleaning database table ${db}..."
    sqlite3 ${databaseFile} "delete from ${db}"
  done
}

main() {
  clean $1
}

main $1
