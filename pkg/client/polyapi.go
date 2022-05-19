package client

import (
	"context"
	"fmt"
	"net/http"

	httpClient "github.com/quanxiang-cloud/cabin/tailormade/client"
	"github.com/quanxiang-cloud/faas/pkg/config"
)

const (
	defaultRootPath    = "/system/app/%s/raw/faas"
	createNamespaceURL = "http://polyapi:9090/api/v1/polyapi/inner/createNamespace"
)

type Client struct {
	c http.Client
}

func New(config *config.Config) *Client {
	return &Client{
		c: httpClient.New(config.InternalNet),
	}
}

type CreateNamespaceReq struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

type CreateNamespaceResp struct {
}

func (c *Client) CreateNamespace(ctx context.Context, appID, group, title string) error {
	ns := fmt.Sprintf(defaultRootPath, appID)
	resp := &CreateNamespaceResp{}
	return httpClient.POST(ctx, &c.c, createNamespaceURL+ns, &CreateNamespaceReq{
		Name:  group,
		Title: title,
	}, resp)
}
