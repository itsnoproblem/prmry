#!/bin/bash

source "${BASH_SOURCE%/*}/functions.sh"
DIR=$(scriptDir)

if ! hash goose 2> /dev/null
then
  echo "Installing goose..."
  go install github.com/pressly/goose/v3/cmd/goose@latest
fi

ENVFILE=$(realpath "$DIR/../.env")
if [ -f "$ENVFILE" ]
then
  . $ENVFILE
fi

goose -s -dir=./migrations mysql "$DB_USER:$DB_PASS@tcp($DB_HOST)/$DB_NAME?parseTime=true&multiStatements=true" "$@"
