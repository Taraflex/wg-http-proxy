package main

import (
	"runtime"
)

const AppName = "{{APPNAME}}"
const Version = "{{GITHUB_REF_NAME}}"
const Arch = runtime.GOOS + "/" + runtime.GOARCH //todo load info from debug.ReadBuildInfo() https://icinga.com/blog/2022/05/25/embedding-git-commit-information-in-go-binaries/
