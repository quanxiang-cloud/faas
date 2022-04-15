package k8s

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/openfunction/apis/core/v1beta1"
	"github.com/openfunction/pkg/client/clientset/versioned"
	ginheader "github.com/quanxiang-cloud/cabin/tailormade/header"
	"gopkg.in/yaml.v3"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"strings"
	"time"
)

type Client interface {
	CreateGitToken(ctx context.Context, host, token string) error
	CreateDocker(ctx context.Context, host, username, secret string) error
	Build(ctx context.Context, data *Function) error
}

type client struct {
	client    *kubernetes.Clientset
	ofn       versioned.Interface
	namespace string
}

func NewClient(namespace string) (Client, error) {
	config := ctrl.GetConfigOrDie()
	clientset := kubernetes.NewForConfigOrDie(config)
	ofn := versioned.NewForConfigOrDie(config)
	return &client{
		client:    clientset,
		namespace: namespace,
		ofn:       ofn,
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

type Function struct {
	Version   string
	Host      string
	Project   string
	GroupName string
	Git       *Git
	Docker    *Docker
	Builder   string
	ENV       map[string]string
}

type Docker struct {
	Host      string
	NameSpace string
	Name      string
}

type Git struct {
	Name string
	Host string
}

func (c *client) Build(ctx context.Context, data *Function) error {
	fn := c.ofn.CoreV1beta1().Functions(c.namespace)
	SourceSubPath := "functions/knative/hello-world-go"
	function := &v1beta1.Function{
		ObjectMeta: metav1.ObjectMeta{
			Name: strings.ToLower(data.GroupName) + "-" + data.Project + "-" + data.Version,
		},
		Spec: v1beta1.FunctionSpec{
			Version: &data.Version,
			Image:   data.Docker.Host + data.Docker.NameSpace + strings.ToLower(data.GroupName) + "-" + data.Project + ":" + data.Version,
			ImageCredentials: &v1.LocalObjectReference{
				Name: data.Docker.Name,
			},
			Build: &v1beta1.BuildImpl{
				SuccessfulBuildsHistoryLimit: pointer.Int32Ptr(2),
				FailedBuildsHistoryLimit:     pointer.Int32Ptr(3),
				Timeout: &metav1.Duration{
					Duration: 10 * time.Minute,
				},
				Builder: &data.Builder,
				Env:     data.ENV,
				SrcRepo: &v1beta1.GitRepo{
					Url:           data.Git.Host + data.GroupName + "/" + data.Project + ".git",
					SourceSubPath: &SourceSubPath,
				},
			},
		},
	}
	marshal, err2 := yaml.Marshal(function)
	if err2 != nil {
		panic(err2)
	}
	fmt.Println(string(marshal))
	_, err := fn.Create(ctx, function, metav1.CreateOptions{})
	return err
}
