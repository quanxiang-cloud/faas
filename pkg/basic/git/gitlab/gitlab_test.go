package gitlab

import (
	"context"
	"fmt"
	"testing"

	"github.com/quanxiang-cloud/faas/pkg/basic/git/define"
)

func TestListGITProject(t *testing.T) {
	git, err := New(&define.Config{
		Token:   "z7MjnQ8Kx4vLGVwB5tEr",
		BaseURL: "http://192.168.201.3:30956/",
	})
	if err != nil {
		panic(err)
	}

	p, err := git.GetGroupProjects(context.Background(), 136)
	if err != nil {
		panic(err)
	}

	for _, v := range p {
		fmt.Println(v.ID)
	}
}
