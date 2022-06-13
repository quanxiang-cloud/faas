package k8s

import (
	"context"
	"github.com/quanxiang-cloud/faas/pkg/config"
	"testing"
)

func TestCreatGit(t *testing.T) {
	newConfig, err := config.NewConfig("../../../configs/config.yml")
	if err != nil {
		panic(err)
	}
	c := NewClient(newConfig, "faas")
	ctx := context.Background()
	ctx = context.WithValue(ctx, "Tenant-Id", "justtest")

	var host = "http://192.168.208.3:10022"
	var keyGenHost = "[192.168.208.3]:10022 ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDEuuYInZAnp5BFjUv1gSTjxD2nwbvslUIFFkfvMtMEK5qTNsDTPjVz9UJfCJWgwjcqJdEvVaVMfZL7wVRZT2rky2CRRO6qVtaOyjDOzWK0aoWnsn5/eEvUE3NRljjZOH+cPdHsbInSIeoQESIA3JDVYs2IC2wYG4+2UrCwsOHpHYepues+6jxQtJOkiBoy+s9DX1Fg4eDD4uQI2h7k7EaUruFApBvmnSFKsixo18SjW12nhJ+HmWWrX88NCFEs+2li4VHJEobpkZzcOMleTAWpU7PPq4DXNsLiLdH8L+Z/B9tXf6+IPrwIWjHYhYKmhl5vtEuG1Ms9M+wQRjMKhBtb\n[192.168.208.3]:10022 ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBK7xuRcLe7nBpeZTsBy+V+7gSftTqF3WztTayAsawhzADSviesm69fNIpO6AHptLFuexTku6ZrOYiYVpvkTBLuI=\n[192.168.208.3]:10022 ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOSLwxuzmsnCdfpCh4aniFBT1vSN4QKTqAYbXUtrG7I8"
	var ssh = "-----BEGIN OPENSSH PRIVATE KEY-----\n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\n-----END OPENSSH PRIVATE KEY-----\n"

	err = c.CreateGitSSH(ctx, host, keyGenHost, ssh)
	if err != nil {
		panic(err)
	}
}

//func TestCreatDocker(t *testing.T) {
//	c := NewClient("faas")
//
//	ctx := context.Background()
//	ctx = context.WithValue(ctx, "Tenant-Id", "vvlgo")
//
//	err := c.CreateDocker(ctx, "http://vvlgo.com", "root", "vvlgo")
//	if err != nil {
//		panic(err)
//	}
//}

// func TestBuilder(t *testing.T) {
// 	c := NewClient("builder")

// 	env := make(map[string]string)
// 	env["FUNC_NAME"] = "HelloWorld"
// 	env["FUNC_CLEAR_SOURCE"] = "true"
// 	env["GOPROXY"] = "https://goproxy.cn,direct"
// 	ctx := context.Background()
// 	ctx = context.WithValue(ctx, "Tenant-Id", "qxcr")
// 	function := &Function{
// 		ID:        "1",
// 		Version:   "v007",
// 		Project:   "faasdemo",
// 		GroupName: "root",
// 		Git: &Git{
// 			Name: "demo",
// 			Host: "ssh://git@192.168.201.3:30679/",
// 		},
// 		Docker: &Docker{
// 			NameSpace: "privitetest/",
// 			Name:      "qxcr",
// 			Host:      "qxcr.xyz/",
// 		},
// 		Builder: "openfunction/builder-go:latest",
// 		ENV:     env,
// 	}
// 	err := c.Build(ctx, function)
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func TestDelFunction(t *testing.T) {
// 	c := NewClient("builder")
// 	ctx := context.Background()
// 	ctx = context.WithValue(ctx, "Tenant-Id", "qxcr")

// 	err := c.DelFunction(ctx, &DelFunction{
// 		Name: "openfunction-samples-v220",
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// }

// // qxcr.xyz/faas/sample-go-func:latest
// func TestCreateServing(t *testing.T) {
// 	c := NewClient("serving")

// 	ctx := context.Background()
// 	env := make(map[string]string)
// 	env["FUNC_NAME"] = "HelloWorld"
// 	env["FUNC_CLEAR_SOURCE"] = "true"
// 	function := &Function{
// 		ID:        "1",
// 		Version:   "v1",
// 		Project:   "demo",
// 		GroupName: "",
// 		Docker: &Docker{
// 			NameSpace: "faas/",
// 			Name:      "qxcr",
// 			Host:      "qxcr.xyz/",
// 		},
// 		ENV: env,
// 	}
// 	err := c.CreateServing(ctx, function)
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func TestDelServing(t *testing.T) {
// 	c := NewClient("serving")

// 	ctx := context.Background()
// 	env := make(map[string]string)
// 	env["FUNC_NAME"] = "HelloWorld"
// 	env["FUNC_CLEAR_SOURCE"] = "true"
// 	function := &Function{
// 		ID:        "1",
// 		Version:   "v207",
// 		Project:   "samples",
// 		GroupName: "OpenFunction",
// 		Docker: &Docker{
// 			NameSpace: "faas/",
// 			Name:      "qxcr",
// 			Host:      "qxcr.xyz/",
// 		},
// 		ENV: env,
// 	}
// 	err := c.DelServing(ctx, function)
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func TestRegister(t *testing.T) {
// 	c := NewClient("builder")

// 	ctx := context.Background()

// 	function := &Function{
// 		Version:   "v1",
// 		Project:   "test",
// 		GroupName: "ifclouddemo",
// 		Git: &Git{
// 			Name: "demo",
// 			Host: "ssh://git@192.168.201.3:30679",
// 		},
// 	}

// 	err := c.RegistAPI(ctx, function, "test")
// 	if err != nil {
// 		panic(err)
// 	}
// }
