package k8s

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	ginheader "github.com/quanxiang-cloud/cabin/tailormade/header"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

type Client interface {
	CreateGitToken(ctx context.Context, host, token string) error
	CreateDocker(ctx context.Context, host, username, secret string) error
}

type client struct {
	client    *kubernetes.Clientset
	namespace string
}

func NewClient(namespace string) (Client, error) {
	config := ctrl.GetConfigOrDie()
	clientset := kubernetes.NewForConfigOrDie(config)
	return &client{
		clientset,
		namespace,
	}, nil
}

const UnexpectedType = "unexpected type"
const Default = "faas"

func (c *client) CreateGitToken(ctx context.Context, host, token string) error {
	secret := c.client.CoreV1().Secrets(c.namespace)
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	if tenantID == "" || tenantID == UnexpectedType {
		tenantID = Default
	}
	data := make(map[string][]byte)
	data["host"] = []byte(host)
	data["token"] = []byte(token)
	s := &v1.Secret{
		Type: v1.SecretTypeOpaque,
		ObjectMeta: ctrl.ObjectMeta{
			Name:      tenantID + "-git",
			Namespace: c.namespace,
		},
		Data: data,
	}
	options := metav1.CreateOptions{}
	_, err := secret.Create(ctx, s, options)
	if err != nil {
		return err
	}
	return err
}

func (c *client) CreateDocker(ctx context.Context, host, username, secret string) error {
	sc := c.client.CoreV1().Secrets(c.namespace)
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	if tenantID == "" || tenantID == UnexpectedType {
		tenantID = Default
	}

	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, secret)))

	val := map[string]map[string]struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
		AUTH     string `json:"auth"`
	}{
		"auths": {
			host: {
				Username: username,
				Password: secret,
				Email:    "",
				AUTH:     auth,
			},
		},
	}
	marshal, err := json.Marshal(val)
	if err != nil {
		return err
	}
	data := make(map[string][]byte)
	data[".dockerconfigjson"] = marshal

	s := &v1.Secret{
		Type: v1.SecretTypeDockerConfigJson,
		ObjectMeta: ctrl.ObjectMeta{
			Name:      tenantID + "-docker",
			Namespace: c.namespace,
		},
		Data: data,
	}
	options := metav1.CreateOptions{}
	_, err = sc.Create(ctx, s, options)
	if err != nil {
		return err
	}
	return err
}
