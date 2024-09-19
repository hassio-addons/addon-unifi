package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	b, err := newBuild()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(3)
	}

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.Usage = func() { b.flagUsage(fs, "") }
	b.init(fs)

	_ = fs.Parse(os.Args[1:])

	if err := b.Validate(); err != nil {
		b.flagUsage(fs, "ERROR: %v", err)
		os.Exit(3)
	}

	err = func() error {
		args := fs.Args()
		if len(args) < 1 {
			return flagFail("command missing")
		}

		cmd, ok := b.commands[args[0]]
		if !ok {
			return flagFail("unknown command: %s", args[0])
		}

		ctx, cancelFunc := context.WithCancel(context.Background())
		sigs := make(chan os.Signal, 1)
		defer close(sigs)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sigs
			cancelFunc()
		}()

		return cmd.f(ctx, args[1:])
	}()

	var fe *flagError
	if errors.As(err, &fe) {
		b.flagUsage(fs, "ERROR: %v", err)
		os.Exit(2)
	} else if err != nil {
		w := flag.CommandLine.Output()
		_, _ = fmt.Fprintf(w, "%s: ERROR %v\n", os.Args[0], err)
		os.Exit(3)
	}
}

func flagFail(msg string, args ...any) *flagError {
	err := flagError(fmt.Sprintf(msg, args...))
	return &err
}

type stringSlice []string

func (this stringSlice) String() string {
	return strings.Join(this, ",")
}

func (this *stringSlice) Set(s string) error {
	var buf stringSlice
	for _, plain := range strings.Split(s, ",") {
		plain = strings.TrimSpace(plain)
		if plain != "" {
			buf = append(buf, plain)
		}
	}
	*this = buf
	return nil
}

type flagError string

func (this flagError) Error() string {
	return string(this)
}
