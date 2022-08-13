#!/bin/bash

set -e

APPNAME=$(basename "$(grep module -i go.mod)")
echo "APPNAME=$APPNAME" >> $GITHUB_ENV

envs=( "android/arm64" "darwin/amd64" "darwin/arm64" "freebsd/386" "freebsd/amd64" "freebsd/arm" "freebsd/arm64" "linux/386" "linux/amd64" "linux/arm" "linux/arm64" "linux/mips" "linux/mips64" "linux/mips64le" "linux/mipsle" "linux/ppc64" "linux/ppc64le" "linux/riscv64" "linux/s390x" "openbsd/amd64" "openbsd/arm64" "openbsd/mips64" "windows/386" "windows/amd64" "windows/arm" "windows/arm64")

for env in "${envs[@]}"
do
    echo $env
    GOOS=$(dirname $env)
    GOARCH=$(basename $env)
    if [[ $GOOS == "windows" ]]; then
        SUFFIX=".exe"
    fi
    
    go build -ldflags="-w -s" -o $APPNAME$SUFFIX . && tar -czf $APPNAME-$GITHUB_REF_NAME-$GOOS-$GOARCH.tar.gz $APPNAME$SUFFIX README.md LICENSE
done