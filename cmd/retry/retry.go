package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/Songmu/retry"
	"github.com/Songmu/wrapcommander"
)

const version = "0.0.1"

type opts struct {
	retry    uint
	interval float64
	cmd      []string
}

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	o, err := parseArgs(args)
	if err != nil {
		if err.Error() != "" {
			fmt.Fprintln(os.Stderr, err)
		}
		writeHelp()
		return 2
	}
	err = o.run()
	return wrapcommander.ResolveExitCode(err)
}

func writeHelp() {
	fmt.Fprintf(os.Stderr, `Usage:
    $ retry COUNT [INTERVAL] COMMAND [ARG]...

Start COMMAND, and retry it until success up to COUNT with INTERVAL seconds.

    COUNT: max retry count
    INTERVAL: interval seconds (default 1.0)

Example:
    $ retry 3 check-hoge ...

Version: %s
`, version)
}

func parseArgs(args []string) (*opts, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("")
	}
	o := &opts{}
	retry, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("retry COUNT should be an integer")
	}
	o.retry = uint(retry)
	cmdIdx := 1
	interval, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		o.interval = 1.0 // default
	} else {
		// custom interval
		o.interval = interval
		cmdIdx++
	}
	o.cmd = args[cmdIdx:]
	if len(o.cmd) < 1 {
		return nil, fmt.Errorf("no command arguments are specified")
	}
	return o, nil
}

func (o *opts) run() error {
	return retry.Retry(uint(o.retry), time.Duration(o.interval*float64(time.Second)), func() error {
		cmd := exec.Command(o.cmd[0], o.cmd[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err := cmd.Run()
		if !wrapcommander.IsInvoked(err) {
			fmt.Fprintln(os.Stderr, err)
		}
		return err
	})
}
