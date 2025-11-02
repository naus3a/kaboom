package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

const Version = "0.0.1"

const ExtShare = ".shab"
const ExtKey = ".keyb"
const ExtPlain = ".plab"
const ExtSignKey = ".sigb"

type AllowedArgType interface {
	bool | uint | string
}

func InitCli(usage string) {
	log.SetFlags(0)
	flag.Usage = func() { fmt.Fprintf(os.Stderr, "%s\n", usage) }
}

func AddArg[T AllowedArgType](arg *T, defaultValue T, commands ...string) error {
	switch t := any(arg).(type) {
	case *bool:
		for i := 0; i < len(commands); i++ {
			flag.BoolVar(t, commands[i], any(defaultValue).(bool), "")
		}
		return nil
	case *uint:
		for i := 0; i < len(commands); i++ {
			flag.UintVar(t, commands[i], any(defaultValue).(uint), "")
		}
		return nil
	case *string:
		for i := 0; i < len(commands); i++ {
			flag.StringVar(t, commands[i], any(defaultValue).(string), "")
		}
		return nil
	default:
		return fmt.Errorf("unsupported arg type")
	}
}

func ReportErrorAndExit(err error) {
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
}

func UnpackCsvArg(arg *string) ([]string, error) {
	vals := strings.Split(*arg, ",")
	if len(vals) < 1 {
		return nil, fmt.Errorf("no values specified")
	}
	return vals, nil
}
