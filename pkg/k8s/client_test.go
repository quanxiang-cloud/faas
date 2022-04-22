package k8s

import (
	"context"
	"fmt"
	"testing"
)

func TestCreatGit(t *testing.T) {
	c := NewClient("faas")
	ctx := context.Background()
	ctx = context.WithValue(ctx, "Tenant-Id", "vvlgo")

	err := c.CreateGitToken(ctx, "http://vvlgo.com", "123")
	if err != nil {
		panic(err)
	}
}

func TestCreatDocker(t *testing.T) {
	c := NewClient("faas")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "Tenant-Id", "vvlgo")

	err := c.CreateDocker(ctx, "http://vvlgo.com", "root", "vvlgo")
	if err != nil {
		panic(err)
	}
}

func TestBuilder(t *testing.T) {
	c := NewClient("builder")

	env := make(map[string]string)
	env["FUNC_NAME"] = "HelloWorld"
	env["FUNC_CLEAR_SOURCE"] = "true"
	env["GOPROXY"] = "https://goproxy.cn,direct"
	ctx := context.Background()
	ctx = context.WithValue(ctx, "Tenant-Id", "qxcr")
	function := &Function{
		ID:        "1",
		Version:   "v007",
		Project:   "faasdemo",
		GroupName: "root",
		Git: &Git{
			Name: "demo",
			Host: "ssh://git@192.168.201.3:30679/",
		},
		Docker: &Docker{
			NameSpace: "privitetest/",
			Name:      "qxcr",
			Host:      "qxcr.xyz/",
		},
		Builder: "openfunction/builder-go:latest",
		ENV:     env,
	}
	err := c.Build(ctx, function)
	if err != nil {
		panic(err)
	}
}

func TestDelFunction(t *testing.T) {
	c := NewClient("builder")
	ctx := context.Background()
	ctx = context.WithValue(ctx, "Tenant-Id", "qxcr")

	err := c.DelFunction(ctx, &DelFunction{
		Name: "openfunction-samples-v220",
	})
	if err != nil {
		panic(err)
	}
}

func TestCreateServing(t *testing.T) {
	c := NewClient("serving")

	ctx := context.Background()
	env := make(map[string]string)
	env["FUNC_NAME"] = "HelloWorld"
	env["FUNC_CLEAR_SOURCE"] = "true"
	function := &Function{
		ID:        "1",
		Version:   "v207",
		Project:   "samples",
		GroupName: "OpenFunction",
		Docker: &Docker{
			NameSpace: "faas/",
			Name:      "qxcr",
			Host:      "qxcr.xyz/",
		},
		ENV: env,
	}
	err := c.CreateServing(ctx, function)
	if err != nil {
		panic(err)
	}
}

func TestDelServing(t *testing.T) {
	c := NewClient("serving")

	ctx := context.Background()
	env := make(map[string]string)
	env["FUNC_NAME"] = "HelloWorld"
	env["FUNC_CLEAR_SOURCE"] = "true"
	function := &Function{
		ID:        "1",
		Version:   "v207",
		Project:   "samples",
		GroupName: "OpenFunction",
		Docker: &Docker{
			NameSpace: "faas/",
			Name:      "qxcr",
			Host:      "qxcr.xyz/",
		},
		ENV: env,
	}
	err := c.DelServing(ctx, function)
	if err != nil {
		panic(err)
	}
}

func TestRegister(t *testing.T) {
	c := NewClient("builder")

	ctx := context.Background()

	function := &Function{
		Version:   "v207",
		Project:   "samples",
		GroupName: "OpenFunction",
		Git: &Git{
			Name: "test",
			Host: "https://github.com",
		},
	}

	err := c.RegistAPI(ctx, function)
	if err != nil {
		panic(err)
	}
}

func TestReverse(t *testing.T) {
	testCase := []string{
		"a-b-c",
		"a-b-c-",
		"-a-b-c",
		"a-b",
	}

	for index, val := range testCase {
		ret, err := reverseName(val)
		if err != nil {
			panic(err)
		}
		fmt.Printf("case %d: %s\n", index, ret)
	}
}
