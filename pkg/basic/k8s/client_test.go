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
	c := NewClient(newConfig)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "Tenant-Id", "justtest")

	var host = "http://192.168.208.3:10022"
	var keyGenHost = "[192.168.208.3]:10022 ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDEuuYInZAnp5BFjUv1gSTjxD2nwbvslUIFFkfvMtMEK5qTNsDTPjVz9UJfCJWgwjcqJdEvVaVMfZL7wVRZT2rky2CRRO6qVtaOyjDOzWK0aoWnsn5/eEvUE3NRljjZOH+cPdHsbInSIeoQESIA3JDVYs2IC2wYG4+2UrCwsOHpHYepues+6jxQtJOkiBoy+s9DX1Fg4eDD4uQI2h7k7EaUruFApBvmnSFKsixo18SjW12nhJ+HmWWrX88NCFEs+2li4VHJEobpkZzcOMleTAWpU7PPq4DXNsLiLdH8L+Z/B9tXf6+IPrwIWjHYhYKmhl5vtEuG1Ms9M+wQRjMKhBtb\n[192.168.208.3]:10022 ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBK7xuRcLe7nBpeZTsBy+V+7gSftTqF3WztTayAsawhzADSviesm69fNIpO6AHptLFuexTku6ZrOYiYVpvkTBLuI=\n[192.168.208.3]:10022 ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOSLwxuzmsnCdfpCh4aniFBT1vSN4QKTqAYbXUtrG7I8\n"
	var ssh = "-----BEGIN OPENSSH PRIVATE KEY-----\nb3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAABFwAAAAdzc2gtcn\nNhAAAAAwEAAQAAAQEA7jOgHwIfF1+GLxiiVNbVA3n2pHc6mXNs0rExHMsB4yxH8UGyBoEq\nGhMJsftqWfrQMozSXVaxEoV6EjxdZWJ+NeTEppqNL5BdgNzkvBYmXEKil0YNi39ylSkean\nZUd3B/LTNMxcgMfSh4DJyBS46FSpjw46JwjPX249T9j6/RNO7r57jtHCzm5uy53h+A8Bcp\nxqxRPF0DcKNKArX8Fn9+f+NQSyN74OVSLASF9fdVVmKuIpP766QkGBREp+EUH5dZ1qVuju\nAG72ypO5eTZv28zOUGI+2eUBWukxRQjx2X6DyYuLiuvw0SQ5YwuOtdMlmjN+8R3UZOCBO5\nnlsjC9IXgQAAA9hWvRQLVr0UCwAAAAdzc2gtcnNhAAABAQDuM6AfAh8XX4YvGKJU1tUDef\nakdzqZc2zSsTEcywHjLEfxQbIGgSoaEwmx+2pZ+tAyjNJdVrEShXoSPF1lYn415MSmmo0v\nkF2A3OS8FiZcQqKXRg2Lf3KVKR5qdlR3cH8tM0zFyAx9KHgMnIFLjoVKmPDjonCM9fbj1P\n2Pr9E07uvnuO0cLObm7LneH4DwFynGrFE8XQNwo0oCtfwWf35/41BLI3vg5VIsBIX191VW\nYq4ik/vrpCQYFESn4RQfl1nWpW6O4AbvbKk7l5Nm/bzM5QYj7Z5QFa6TFFCPHZfoPJi4uK\n6/DRJDljC4610yWaM37xHdRk4IE7meWyML0heBAAAAAwEAAQAAAQEAksirVhBXmm6Z4HG8\nrq0H7LF+hFCYgrD3EUpxaj3K9vN1jbye+JUlnZW3tr6oBbSiEVRn2W6vhStNUulx4CT2gz\n1q1QOOFw5OIDd+mEzJ7LZ/PNrFQ/4QXMxO2un6GPPw2qivGX2x/HowqAzVO/siNhrS7mNO\nGf8D2deJPL7qG0QjB1axasUYzk5rRQyPg3zWl0aKDcbqfAzFt4L32jGw1jV2w61Jm92LiL\n8FzNz8oX/16sk/HC4W/Ps6L+x6h82jP0nZj2QOyp0K3zi6+peIz1FAeJS8wbuu8/uUpXeu\n0EMpgQb0eiHZyXldW89ddzEhPtVBhN1cvy54aCiNpMF3VQAAAIAXmdwFRB2nt4j5Qq4ZrB\nIzB1Kxj7lENscWGIfXNTj2XQj3h3tlzEp0vpkPeTCHTetpBJitWpnRwVDo1s609ZYuJQ4f\n2Nn6VIS6VJejXWaXYKvC0qThCxMhFW7GMZwuSlvVgHQhELsS5C1zBL5NsHorkMIupIWfzk\nP3BLML5vtc+gAAAIEA/yo/JDbmdu9Cvx6nVlTF28XzcLr8EvN68lOF9d8aOixV3SvLl1i5\neOQrHx0K0NvuTPD/og8xrCjvX9kfJRWggZEYiKtB7QCRTrVDdkLc5S6miuJFRtmNfqsLR6\n0ZAjOsa0Uh/NTuaP3nR1a0hB0UkQkYVv+arb9vLcYMomm5TmsAAACBAO77KyNzkckHFDee\nRfBm3diMCISWi3QOBvTgcYqwAQjeCDrOVKuyEi3+N+ZgNj/cr6cjwMz31I6R+okc/TaMGK\nM73SGrrA8AF/bEMxlfrGlaNrLTOYpUlXtjBvu0LbkYow54wIQ0GznbKXVAlnHKxs7tMgWV\nMPGdqO3auNi5YRTDAAAAHHZ2bGdvQFZWTEdPLU1hY2Jvb2tQcm8ubG9jYWwBAgMEBQY=\n-----END OPENSSH PRIVATE KEY-----\n"
	err = c.CreateGitSSH(ctx, host, keyGenHost, ssh)
	if err != nil {
		panic(err)
	}
}

func TestCreatDocker(t *testing.T) {
	newConfig, err := config.NewConfig("../../../configs/config.yml")
	if err != nil {
		panic(err)
	}
	c := NewClient(newConfig)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "Tenant-Id", "justdocker")

	err = c.CreateDocker(ctx, "http://qxcr.xyz", "qxptest", "ZHU**jie9")
	if err != nil {
		panic(err)
	}
}

func TestBuilder(t *testing.T) {
	newConfig, err := config.NewConfig("../../../configs/config.yml")
	if err != nil {
		panic(err)
	}
	c := NewClient(newConfig)

	env := make(map[string]string)
	env["FUNC_NAME"] = "HelloWorld"
	env["FUNC_CLEAR_SOURCE"] = "true"
	env["GOPROXY"] = "https://goproxy.cn,direct"
	ctx := context.Background()
	ctx = context.WithValue(ctx, "Tenant-Id", "justtest")
	function := &Function{
		ID:        "1",
		Version:   "v007",
		Project:   "faasdemo",
		GroupName: "root",
		Git: &Git{
			Name: "demo",
			Host: "ssh://git@192.168.208.3:10022/",
		},
		Docker: &Docker{
			NameSpace: "testcloud/",
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
