#!/bin/bash

goose -s -dir=./migrations mysql "root:root@tcp(127.0.0.1)/rgb?parseTime=true&multiStatements=true" "$@"
