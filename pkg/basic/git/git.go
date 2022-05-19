package git

import (
	"fmt"

	"github.com/quanxiang-cloud/faas/pkg/basic/git/define"
	"github.com/quanxiang-cloud/faas/pkg/basic/git/gitlab"
)

type GitType int

const (
	Gitlab GitType = iota
)

func GetClient(gitType GitType, token, baseURL string) (git GIT, err error) {
	switch gitType {
	case Gitlab:
		git, err = gitlab.New(&define.Config{
			Token:   token,
			BaseURL: baseURL,
		})
	default:
		err = fmt.Errorf("not support type(%d)", gitType)
	}
	return
}
