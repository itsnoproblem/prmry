#!/bin/bash

source "${BASH_SOURCE%/*}/functions.sh"
DIR=$(scriptDir)

if ! hash templ 2> /dev/null
then
  echo "Installing templ..."
  go install github.com/a-h/templ/cmd/templ@latest
fi

ulimit -n 99999
cd "$DIR"/../ || (echo "failed to change directory to $DIR" && exit)
templ generate "$DIR"/../ && go run "$DIR"/../cmd/http-server
cd -
