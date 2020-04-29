#!/usr/bin/env bash

# Run all the tests and generate coverage file.

set -e
echo "" > coverage.txt

if [ "$1" ]; then
    echo "running tests on circleci"
    export PROFILE="ci"
fi


for d in $(go list ./... | grep -v vendor); do
    go test -race -coverprofile=profile.out -covermode=atomic $d
    if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm profile.out
    fi
done

if [ "$2" ]; then
    echo "Publishing test coverage"
    cd ../ && bash <(curl -s https://codecov.io/bash) -t $2
fi

