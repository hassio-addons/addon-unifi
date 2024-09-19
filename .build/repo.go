package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func newRepo() (repo, error) {
	result := repo{
		packages: newPackages(),
		prs:      newPrs(),
		actions:  newActions(),
	}
	if v, ok := os.LookupEnv("GITHUB_OWNER_TYPE"); ok {
		if err := result.ownerType.Set(v); err != nil {
			return repo{}, fmt.Errorf("GITHUB_OWNER_TYPE: %w", err)
		}
	}
	if v, ok := os.LookupEnv("GITHUB_REPOSITORY"); ok {
		parts := strings.Split(v, "/")
		if len(parts) != 2 {
			return repo{}, fmt.Errorf("GITHUB_REPOSITORY: illegal github repository: %s", v)
		}
		if err := result.owner.Set(parts[0]); err != nil {
			return repo{}, fmt.Errorf("GITHUB_REPOSITORY: %w", err)
		}
		if err := result.name.Set(parts[1]); err != nil {
			return repo{}, fmt.Errorf("GITHUB_REPOSITORY: %w", err)
		}
	}
	if v, ok := os.LookupEnv("GITHUB_REPOSITORY_OWNER"); ok {
		if err := result.owner.Set(v); err != nil {
			return repo{}, fmt.Errorf("GITHUB_REPOSITORY_OWNER: %w", err)
		}
	}
	if v, ok := os.LookupEnv("GITHUB_REPO"); ok {
		if err := result.name.Set(v); err != nil {
			return repo{}, fmt.Errorf("GITHUB_REPO: %w", err)
		}
	}
	return result, nil
}

func (this *repo) init(b *build, fs *flag.FlagSet) {
	this.build = b
	fs.Var(&this.ownerType, "ownerType", "Can be either 'user' or 'org'")
	fs.Var(&this.owner, "owner", "")
	fs.Var(&this.name, "repo", "")
	this.packages.init(b, fs)
	this.prs.init(b, fs)
	this.actions.init(b, fs)
}

func (this *repo) Validate() error {
	if err := this.ownerType.Validate(); err != nil {
		return err
	}
	if err := this.owner.Validate(); err != nil {
		return err
	}
	if err := this.name.Validate(); err != nil {
		return err
	}
	if err := this.packages.Validate(); err != nil {
		return err
	}
	if err := this.prs.Validate(); err != nil {
		return err
	}
	if err := this.actions.Validate(); err != nil {
		return err
	}
	return nil
}

type repo struct {
	build *build

	ownerType ownerType
	owner     owner
	name      repoName

	packages packages
	prs      prs
	actions  actions
}

func (this repo) String() string {
	return fmt.Sprintf("%v:%s/%s", this.ownerType, this.owner, this.name)
}

func (this repo) SubName(sub string) string {
	if sub == "" {
		return this.name.String()
	}
	return fmt.Sprintf("%v%%2F%s", this.name, sub)
}

func (this repo) SubString(sub string) string {
	if sub == "" {
		return this.String()
	}
	return fmt.Sprintf("%v/%s", this, sub)
}

type ownerType uint8

func (this ownerType) String() string {
	switch this {
	case user:
		return "user"
	case org:
		return "org"
	default:
		return fmt.Sprintf("illegal-owner-type-%d", this)
	}
}

func (this ownerType) Validate() error {
	switch this {
	case user:
		return nil
	case org:
		return nil
	default:
		return fmt.Errorf("illegal-owner-type-%d", this)
	}
}

func (this *ownerType) Set(v string) error {
	switch v {
	case "user":
		*this = user
	case "org":
		*this = org
	default:
		return fmt.Errorf("unknown ownerType: %s", v)
	}
	return nil
}

const (
	user ownerType = iota
	org
)

type owner string

func (this owner) String() string {
	return string(this)
}

var ownerRegex = regexp.MustCompile("^[a-zA-Z0-9](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?$")

func (this *owner) Set(v string) error {
	buf := owner(v)
	if err := buf.Validate(); err != nil {
		return err
	}
	*this = buf
	return nil
}

func (this owner) Validate() error {
	if this == "" {
		return fmt.Errorf("no owner provided")
	}
	if !ownerRegex.MatchString(string(this)) {
		return fmt.Errorf("illegal owner: %s", this)
	}
	return nil
}

type repoName string

func (this repoName) String() string {
	return string(this)
}

var repoNameRegex = regexp.MustCompile("^[a-zA-Z0-9-_.]+$")

func (this *repoName) Set(v string) error {
	buf := repoName(v)
	if err := buf.Validate(); err != nil {
		return err
	}
	*this = repoName(v)
	return nil
}

func (this repoName) Validate() error {
	if this == "" {
		return fmt.Errorf("no repo name provided")
	}
	if !repoNameRegex.MatchString(string(this)) {
		return fmt.Errorf("illegal repo name: %s", this)
	}
	return nil
}
