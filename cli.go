package main

import (
	"errors"
	"flag"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var LogLevel = struct {
	FATAL uint64
	ERROR uint64
	INFO  uint64
	DEBUG uint64
}{
	FATAL: 0b1,
	ERROR: 0b11,
	INFO:  0b1111,
	DEBUG: 0b11111,
}

type Cli struct {
	ConfigFile   string
	Port         uint64
	LogLevel     uint64
	StartAndExit bool
}

func (f *Cli) SPort() string {
	return strconv.FormatUint(f.Port, 10)
}

func ParseFlags() (Cli, error) {

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s %s %s\nUsage:\n", AppName, Version, Arch)
		flag.PrintDefaults()
	}

	var cli Cli

	flag.Uint64Var(&cli.Port, "p", 41970, "Proxy port [0-49151]")
	var ll string
	flag.StringVar(&ll, "l", "info", "LogLevel: [debug, info, error, fatal]")
	flag.BoolVar(&cli.StartAndExit, "x", false, "exit app after proxy server started")
	flag.Parse()

	if len(flag.Args()) < 1 {
		return cli, errors.New("required one non-flag cli arguments for config file location")
	}
	cli.ConfigFile = flag.Args()[0]

	if cli.Port > 49151 {
		return cli, errors.New("port must be <= 49151")
	}
	v := reflect.ValueOf(LogLevel).FieldByName(strings.ToUpper(ll))
	if v.Kind() == 0 {
		return cli, errors.New("logLevel must be one of [debug, info, error, fatal]")
	}
	cli.LogLevel = v.Uint()
	return cli, nil
}
