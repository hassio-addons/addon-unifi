package main

import (
	"context"
	"flag"
	"fmt"
	"iter"
	"log"
	"time"

	"github.com/google/go-github/v65/github"
)

func newActions() actions {
	return actions{}
}

func (this *actions) init(b *build, _ *flag.FlagSet) {
	this.build = b
}

type actions struct {
	build *build
}

func (this *actions) Validate() error { return nil }

func (this *actions) workflowByFilename(ctx context.Context, fn string) (*workflow, error) {
	v, _, err := this.build.client.Actions.GetWorkflowByFileName(ctx, this.build.owner.String(), this.build.name.String(), fn)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve workflow %s from %v: %w", fn, this.build, err)
	}
	return &workflow{
		v,
		this,
	}, nil
}

type workflow struct {
	*github.Workflow

	parent *actions
}

func (this *workflow) String() string {
	return fmt.Sprintf("%s(%d)@%v", *this.Name, *this.ID, this.parent.build.repo)
}

func (this *workflow) runs(ctx context.Context) iter.Seq2[*workflowRun, error] {
	return func(yield func(*workflowRun, error) bool) {
		var opts github.ListWorkflowRunsOptions
		opts.PerPage = 100

		for {
			candidates, rsp, err := this.parent.build.client.Actions.ListWorkflowRunsByID(ctx, this.parent.build.owner.String(), this.parent.build.name.String(), *this.ID, &opts)
			if err != nil {
				yield(nil, fmt.Errorf("cannot retrieve workflow runs of %v (page: %d): %w", this, opts.Page, err))
				return
			}
			for _, v := range candidates.WorkflowRuns {
				if !yield(&workflowRun{v, this}, nil) {
					return
				}
			}
			if rsp.NextPage == 0 {
				return
			}
			opts.Page = rsp.NextPage
		}
	}
}

type workflowRun struct {
	*github.WorkflowRun
	parent *workflow
}

func (this *workflowRun) reload(ctx context.Context) error {
	v, _, err := this.parent.parent.build.client.Actions.GetWorkflowRunByID(ctx, this.parent.parent.build.owner.String(), this.parent.parent.build.name.String(), *this.ID)
	if err != nil {
		return fmt.Errorf("cannot reload workflow run %v: %w", this, err)
	}
	this.WorkflowRun = v
	return nil
}

func (this *workflowRun) wait(ctx context.Context) error {
	start := time.Now()
	for ctx.Err() == nil {
		if err := this.reload(ctx); err != nil {
			return err
		}
		if this.Status != nil && *this.Status == "completed" {
			return nil
		}
		log.Printf("INFO %v is still running (waiting since %v)...", this, time.Since(start).Truncate(time.Second))
		time.Sleep(this.parent.parent.build.waitTimeout)
	}
	return ctx.Err()
}

func (this *workflowRun) rerun(ctx context.Context) error {
	_, err := this.parent.parent.build.client.Actions.RerunWorkflowByID(ctx, this.parent.parent.build.owner.String(), this.parent.parent.build.name.String(), *this.ID)
	if err != nil {
		return fmt.Errorf("cannot rerun workflow run %v: %w", this, err)
	}
	return nil
}

func (this workflowRun) String() string {
	return fmt.Sprintf("%d@%v", *this.ID, this.parent)
}
