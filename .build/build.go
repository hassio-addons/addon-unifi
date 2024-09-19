package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/google/go-github/v65/github"
)

func newBuild() (build, error) {
	r, err := newRepo()
	if err != nil {
		return build{}, err
	}
	return build{
		repo:        r,
		waitTimeout: time.Second * 3,
		commands:    make(map[string]command),
	}, nil
}

type build struct {
	repo

	waitTimeout time.Duration

	client *github.Client

	commands map[string]command
}

type command struct {
	f     func(context.Context, []string) error
	usage string
}

func (this *build) init(fs *flag.FlagSet) {
	this.client = github.NewClient(nil).
		WithAuthToken(os.Getenv("GITHUB_TOKEN"))

	fs.DurationVar(&this.waitTimeout, "wait-timeout", this.waitTimeout, "")

	this.repo.init(this, fs)
}

func (this *build) Validate() error {
	if err := this.repo.Validate(); err != nil {
		return err
	}
	return nil
}

func (this *build) registerCommand(name, usage string, action func(context.Context, []string) error) {
	this.commands[name] = command{action, usage}
}

func (this *build) flagUsage(fs *flag.FlagSet, reasonMsg string, args ...any) {
	w := fs.Output()
	_, _ = fmt.Fprint(w, `Usage of .build:
`)
	if reasonMsg != "" {
		_, _ = fmt.Fprintf(w, "Error: %s\n", fmt.Sprintf(reasonMsg, args...))
	}
	_, _ = fmt.Fprintf(w, "Syntax: %s [flags] <command> [commandSpecificArgs]\nCommands:\n", fs.Name())
	for n, c := range this.commands {
		_, _ = fmt.Fprintf(w, "\t%s %s\n", n, c.usage)
	}
	_, _ = fmt.Fprint(w, "Flags:\n")
	fs.PrintDefaults()
}
