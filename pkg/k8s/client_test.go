package k8s

import (
	"context"
	"testing"
)

func TestCreatGit(t *testing.T) {

	c, err := NewClient("lowcode")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	err = c.CreateGitToken(ctx, "http://vvlgo.com", "123")
	if err != nil {
		panic(err)
	}
}

func TestCreatDocker(t *testing.T) {
	c, err := NewClient("lowcode")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	err = c.CreateDocker(ctx, "http://vvlgo.com", "root", "vvlgo")
	if err != nil {
		panic(err)
	}
}
