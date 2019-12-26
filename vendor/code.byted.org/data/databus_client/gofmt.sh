#!/bin/bash

gofmt -w -e -l -d *.go

git diff --quiet

if [[ $? -eq 1 ]]; then
    exit 1
else
    exit 0
fi

