#!/bin/bash

source "${BASH_SOURCE%/*}/functions.sh"
DIR=$(scriptDir)

if ! hash reflex 2> /dev/null
then
  echo "Installing reflex..."
  go install github.com/cespare/reflex@latest
fi

cd "$DIR"/../ || (echo "failed to change directory to $DIR" && exit 1)
reflex -d fancy -s -r '\.templ' "$DIR"/server.sh
cd -
