#!/bin/bash

export LANG=en_US.UTF-8

depth_foler() {
    this_dir=`pwd`
    source_folder="$1"
    source_folder=`echo $source_folder |sed 's#/$##g'`
    test_folder="$2"
    cd $test_folder
    count=0
    while [ ! `pwd` = $source_folder ]
    do
        count=`expr $count + 1`
        cd ..
    done
    cd $this_dir
    return $count
}

function getDirMaxDepth() {
    folder_name="$1"
    if [ ! -d "$folder_name" ];then
        echo "The dir(\"${folder_name}\") does not exist!"
        exit 3
    fi
    this_dir_tmp=`pwd`
    folder_name=`echo "$folder_name"|sed "s#^./#$this_dir_tmp/#g"`
    folder_name=`echo "$folder_name"|sed "s#^\([a-zA-Z]\+.*\)#$this_dir_tmp/\1#g"`
    target_folder="$folder_name"
    depth_max=1
    for i in `du "$target_folder"` ;do
        if [ -d $i -a ! $i = $target_folder ];then
            depth_foler "$target_folder" "$i"
            retval=$?
            #echo "$i : $retval"
            if [ $depth_max -lt $retval ];then
                    depth_max=$retval
            fi
        fi
    done
    echo -n $depth_max
}

changeDirCharSet() {
  targetDir=$1
  maxDepth=$2
  for ((i=${maxDepth}; i>=1; i--)); do
    for dir in $(find "${targetDir}" -maxdepth ${i} -mindepth ${i} -type d); do
      dirName=$(dirname $dir)
      baseName=$(basename $dir)
      afterChange=$(echo -n ${baseName} | iconv -f gbk -t utf8)
      if [[ ! ${afterChange} == ${baseName} ]]; then
        pushd ${dirName} &>/dev/null
        mv ${baseName} ${afterChange}
        echo change ${dir} to ${dirName}/${afterChange}
        popd &>/dev/null
      fi
    done
  done
}

changeFileCharSet() {
  targetDir=$1
  maxDepth=$2
  for ((i=$((${maxDepth}+1)); i>=1; i--)); do
    for dir in $(find "${targetDir}" -maxdepth ${i} -mindepth ${i} -type f); do
      dirName=$(dirname $dir)
      baseName=$(basename $dir)
      afterChange=$(echo -n ${baseName} | iconv -f gbk -t utf8)
      if [[ ! ${afterChange} == ${baseName} ]]; then
        pushd ${dirName} &>/dev/null
        mv ${baseName} ${afterChange}
        echo change ${dir} to ${dirName}/${afterChange}
        popd &>/dev/null
      fi
    done
  done
}

main() {
  maxDepth=$(getDirMaxDepth "${1}")
  echo ${maxDepth}
  changeDirCharSet $1 ${maxDepth}
  changeFileCharSet $1 ${maxDepth}
}

main $1
