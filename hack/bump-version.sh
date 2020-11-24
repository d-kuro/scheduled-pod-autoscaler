#!/bin/bash

set -e -o pipefail

if [ $# -ne 2 ]; then
  echo 'please specify two tags, "bump-version.sh <from> <to>"'
  exit 1
fi

# works only with BSD sed
git grep -l -e "$1" --and -e 'd-kuro/scheduled-pod-autoscaler' | xargs -I {} sed -i '' -e "s/$1/$2/g" {}
