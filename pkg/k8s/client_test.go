package k8s

import (
	"context"
	"testing"
)

func TestCreatGit(t *testing.T) {

	c, err := NewClient("faas")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, "Tenant-Id", "vvlgo")

	err = c.CreateGitToken(ctx, "http://vvlgo.com", "123")
	if err != nil {
		panic(err)
	}
}

func TestCreatDocker(t *testing.T) {
	c, err := NewClient("faas")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, "Tenant-Id", "vvlgo")

	err = c.CreateDocker(ctx, "http://vvlgo.com", "root", "vvlgo")
	if err != nil {
		panic(err)
	}
}

func TestBuilder(t *testing.T) {
	c, err := NewClient("faas")
	if err != nil {
		panic(err)
	}
	env := make(map[string]string)
	env["FUNC_NAME"] = "HelloWorld"
	env["FUNC_CLEAR_SOURCE"] = "true"
	env["GOPROXY"] = "https://goproxy.cn,direct"
	ctx := context.Background()
	ctx = context.WithValue(ctx, "Tenant-Id", "qxcr")
	function := &Function{
		ID:        "1",
		Version:   "v207",
		Project:   "samples",
		GroupName: "OpenFunction",
		Git: &Git{
			Name: "",
			Host: "https://github.com/",
		},
		Docker: &Docker{
			NameSpace: "faas/",
			Name:      "qxcr",
			Host:      "qxcr.xyz/",
		},
		Builder: "openfunction/builder-go:latest",
		ENV:     env,
	}
	err = c.Build(ctx, function)
	if err != nil {
		panic(err)
	}
}

func TestDelFunction(t *testing.T) {
	c, err := NewClient("faas")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, "Tenant-Id", "qxcr")

	err = c.DelFunction(ctx, &DelFunction{
		Name: "openfunction-samples-v207",
	})
	if err != nil {
		panic(err)
	}
}
