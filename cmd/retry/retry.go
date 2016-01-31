package main

import (
	"os"
	"os/exec"
	"time"

	"github.com/Songmu/retry"
	"github.com/Songmu/wrapcommander"
	"github.com/jessevdk/go-flags"
)

type opts struct {
	Retry    int64   `short:"r" long:"retry" description:"hoge"`
	Interval float64 `short:"i" long:"interval" description:"hoge"`
}

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	p, o, cmdArgs, err := parseArgs(args)
	if err != nil || len(cmdArgs) < 1 {
		if ferr, ok := err.(*flags.Error); !ok || ferr.Type != flags.ErrHelp {
			p.WriteHelp(os.Stderr)
		}
		return 2
	}
	err = o.run(cmdArgs)
	return wrapcommander.ResolveExitCode(err)
}

func parseArgs(args []string) (*flags.Parser, *opts, []string, error) {
	o := &opts{}
	p := flags.NewParser(o, flags.Default)

	rest, err := p.ParseArgs(args)
	return p, o, rest, err
}

func (o *opts) run(args []string) error {
	return retry.Retry(uint(o.Retry), time.Duration(o.Interval*float64(time.Second)), func() error {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	})
}
