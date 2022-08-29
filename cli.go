package main

import (
	"errors"
	"flag"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func getFields[T any](i *T) string {
	v := reflect.ValueOf(i).Elem().Type()
	n := v.NumField()
	r := make([]string, n)
	for i := 0; i < n; i++ {
		r[i] = strings.ToLower(v.Field(i).Name)
	}
	return strings.Join(r, ", ")
}

var LogLevel = struct {
	FATAL uint64
	ERROR uint64
	INFO  uint64
	DEBUG uint64
}{0b1, 0b11, 0b111, 0b1111}

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
	validLogLevels := getFields(&LogLevel)

	flag.Uint64Var(&cli.Port, "p", 41970, "Proxy port [0-49151]")
	var ll string
	flag.StringVar(&ll, "l", "info", "LogLevel: ["+validLogLevels+"]")
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
		return cli, errors.New("logLevel must be one of [" + validLogLevels + "]")
	}
	cli.LogLevel = v.Uint()
	return cli, nil
}
