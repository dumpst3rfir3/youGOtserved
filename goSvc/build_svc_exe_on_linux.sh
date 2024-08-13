#!/usr/bin/env bash

# Setup ENV
export GOOS=windows

# Version Info
sed "s/youGOtserved/${1:-youGOtserved}/g" versioninfo.template > versioninfo.json
sed "s/youGOtserved/${1:-youGOtserved}/g" main.template > main.go
go generate
# Compile EXE for windows
go build --trimpath --buildvcs=false --ldflags="-s -w" .
rm -f versioninfo.json
