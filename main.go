package main

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"

	"github.com/docopt/docopt-go"
)

var (
	version = "[manual build]"
	usage   = "when " + version + `

Run specified command when given conditions becomes real.

Usage:
  when [options] [--dns <nameserver>]... -- <cmd>...
  when -h | --help
  when --version

Options:
  --pre-start <cmd>            Run specified command before starting command.
  --dns <nameserver>           Run when specified DNS server is available.
  --dns-check-domain <domain>  Well-known domain name to check DNS availability.
                                [default: google.com]
  -i --interval <duration>     Interval of checking conditions in milliseconds.
                                [default: 100]
  -h --help                    Show this screen.
  --version                    Show version.
`
)

func main() {
	args, err := docopt.Parse(usage, nil, true, version, false)
	if err != nil {
		panic(err)
	}

	intervalMs, _ := strconv.Atoi(args["--interval"].(string))

	conditions := []*ConditionDNS{}
	for _, nameserver := range args["--dns"].([]string) {
		conditions = append(conditions, NewConditionDNS(
			nameserver, args["--dns-check-domain"].(string),
		))
	}

	total := len(conditions)
	for i, condition := range conditions {
		for !condition.Ready() {
			log.Printf(
				"[%d/%d] checking condition: DNS %s",
				i+1, total, condition.address,
			)

			err := condition.Check()
			if err != nil {
				log.Printf(
					"[%d/%d] DNS %s is not ready: %s",
					i+1, total, condition.address, err,
				)
				time.Sleep(time.Duration(intervalMs) * time.Millisecond)
			} else {
				log.Printf(
					"[%d/%d] DNS %s is ready",
					i+1, total, condition.address,
				)
			}
		}
	}

	if preStart, ok := args["--pre-start"].(string); ok {
		log.Printf("starting pre-start command: %q", preStart)
		cmd := exec.Command("/bin/sh", "-c", preStart)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}

	cmd := args["<cmd>"].([]string)

	log.Printf("starting %q", cmd)

	err = syscall.Exec(cmd[0], cmd, os.Environ())
	if err != nil {
		log.Fatal(err)
	}
}
