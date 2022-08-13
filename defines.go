package main

import (
	"runtime"
	"runtime/debug"
	"strings"
)

var info, _ = debug.ReadBuildInfo()
var parts = strings.Split(info.Main.Path, "/")

var AppName = parts[len(parts)-1]
var GithubPacUrl = "https://" + strings.ToLower(parts[1]) + ".github.io/" + AppName + "/pac.js.br"
var Version = info.Main.Version

const Arch = runtime.GOOS + "/" + runtime.GOARCH
