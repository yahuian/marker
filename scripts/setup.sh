#!/bin/bash

# set up git hooks
SRC=".github/.githooks/pre-commit"
DST=".git/hooks/pre-commit"
cp $SRC $DST
chmod u+x $DST

# install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.46.2
