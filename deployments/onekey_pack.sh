#!/bin/bash

main() {
  bash make_version.sh
  tar zcf git-translator.tar.gz package version conf.json deploy deploy.sh deploy.README
}

main
