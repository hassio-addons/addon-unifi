package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v65/github"
)

func newPrs() prs {
	return prs{}
}

func (this *prs) init(b *build, _ *flag.FlagSet) {
	this.build = b

	b.registerCommand("rerun-pr-workflow", "<prNumber> <workflowFilename>", this.rerunLatestWorkflowCmd)
	b.registerCommand("has-pr-label", "<prNumber> <label>", this.hasLabelCmd)
	b.registerCommand("is-pr-open", "<prNumber>", this.isOpenCmd)
}

func (this *prs) Validate() error { return nil }

type prs struct {
	build *build
}

func (this *prs) rerunLatestWorkflowCmd(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return flagFail("no pr number provided")
	}
	if len(args) < 2 {
		return flagFail("no workflow filename provided")
	}

	number, err := strconv.Atoi(args[0])
	if err != nil {
		return flagFail("illegal pr number provided: %s", args[0])
	}

	v, err := this.byId(ctx, number)
	if err != nil {
		return err
	}

	start := time.Now()
	for ctx.Err() == nil {
		wfr, err := v.latestWorkflowRun(ctx, args[1])
		if err != nil {
			return err
		}
		if wfr.Status != nil && strings.EqualFold(*wfr.Status, "completed") {
			if err := wfr.rerun(ctx); err != nil {
				return err
			}
			fmt.Printf("INFO successfully rerurn worflow run %v (%s) for pr %d (%s)", wfr, *wfr.HTMLURL, *v.ID, *v.HTMLURL)
			return nil
		}
		log.Printf("INFO latest worflow run %v is still running (waiting since %v)...", this, time.Since(start).Truncate(time.Second))
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(this.build.waitTimeout):
		}
	}

	return ctx.Err()
}

func (this *prs) hasLabelCmd(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return flagFail("no pr number provided")
	}
	if len(args) < 2 {
		return flagFail("no label provided")
	}

	number, err := strconv.Atoi(args[0])
	if err != nil {
		return flagFail("illegal pr number provided: %s", args[0])
	}
	v, err := this.byId(ctx, number)
	if err != nil {
		return err
	}

	if v.hasLabel(args[1]) {
		os.Exit(0)
	} else {
		os.Exit(1)
	}

	return nil
}

func (this *prs) isOpenCmd(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return flagFail("no pr number provided")
	}

	number, err := strconv.Atoi(args[0])
	if err != nil {
		return flagFail("illegal pr number provided: %s", args[0])
	}
	v, err := this.byId(ctx, number)
	if err != nil {
		return err
	}

	if v.isOpen() {
		os.Exit(0)
	} else {
		os.Exit(1)
	}

	return nil
}

func (this *prs) byId(ctx context.Context, number int) (*pr, error) {
	v, _, err := this.build.client.PullRequests.Get(ctx, this.build.owner.String(), this.build.name.String(), number)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve pull request %d from %v: %w", number, this.build, err)
	}
	return &pr{
		v,
		this,
	}, nil
}

type pr struct {
	*github.PullRequest

	parent *prs
}

func (this *pr) String() string {
	return fmt.Sprintf("%d@%v", *this.ID, this.parent.build.repo)
}

func (this *pr) hasLabel(label string) bool {
	return slices.ContainsFunc(this.Labels, func(candidate *github.Label) bool {
		if candidate == nil {
			return false
		}
		return candidate.Name != nil && *candidate.Name == label
	})
}

func (this *pr) isOpen() bool {
	return this.State != nil && strings.EqualFold(*this.State, "open")
}

func (this *pr) latestWorkflowRun(ctx context.Context, workflowFn string) (*workflowRun, error) {
	wf, err := this.parent.build.actions.workflowByFilename(ctx, workflowFn)
	if err != nil {
		return nil, fmt.Errorf("cannot get workflow %s of %v: %w", workflowFn, this, err)
	}
	for candidate, err := range wf.runs(ctx) {
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve workflow runs for pr %v: %w", this, err)
		}
		if slices.ContainsFunc(candidate.PullRequests, func(cpr *github.PullRequest) bool {
			return cpr != nil && cpr.ID != nil && *cpr.ID == *this.ID
		}) {
			return candidate, nil
		}
	}
	return nil, nil
}
