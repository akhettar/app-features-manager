#!/usr/bin/env bash
go get github.com/golang/mock/gomock
go get github.com/golang/mock/mockgen

# create dir where all the mocks end
mkdir -p ../mocks


$GOPATH/bin/mockgen -destination=mocks/mock_repository.go -package=mocks github.com/akhettar/app-features-manager/repository Repository
$GOPATH/bin/mockgen -destination=mocks/mock_unleash.go -package=mocks github.com/akhettar/app-features-manager/features UnleashService
