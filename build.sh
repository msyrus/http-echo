#! /bin/sh
set -e

PROJ="http-echo"
ORG_PATH="github.com/msyrus"
REPO_PATH="${ORG_PATH}/${PROJ}"

if ! [ -x "$(command -v go)" ]; then
    echo "go is not installed"
    exit
fi
if [ -z "${GOPATH}" ]; then
    echo "set GOPATH"
    exit
fi

PATH="${PATH}:${GOPATH}/bin"

# Building go binary
go install -v .
