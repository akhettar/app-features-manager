#!/usr/bin/env bash

set -e

echo "creating swagger docs folder"
rm -rf ../docs && mkdir ../docs

echo "Installing swaggo command line tool"
go get -u github.com/swaggo/swag
cd $GOPATH/src/github.com/swaggo/swag/cmd/swag
go install
echo "swag tool install complete."

echo "Generating swagger document"
cd - && swag initls