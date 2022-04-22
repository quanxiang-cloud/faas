package k8s

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/openfunction/apis/core/v1beta1"
	"github.com/openfunction/pkg/client/clientset/versioned"
	ginheader "github.com/quanxiang-cloud/cabin/tailormade/header"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"strings"
	"time"
)

// Client Client
type Client interface {
	CreateGitToken(ctx context.Context, host, token string) error
	CreateGitSSH(ctx context.Context, host, ssh string) error
	CreateDocker(ctx context.Context, host, username, secret string) error
	Build(ctx context.Context, data *Function) error
	DelFunction(ctx context.Context, data *DelFunction) error
}

type client struct {
	client          *kubernetes.Clientset
	ofn             versioned.Interface
	k8sNamespace    string
	dockerNamespace string
}

// NewClient NewClient
func NewClient(namespace string) (Client, error) {
	config := ctrl.GetConfigOrDie()
	clientset := kubernetes.NewForConfigOrDie(config)
	ofn := versioned.NewForConfigOrDie(config)
	return &client{
		client:       clientset,
		k8sNamespace: namespace,
		ofn:          ofn,
	}, nil
}

// UnexpectedType UnexpectedType
const UnexpectedType = "unexpected type"

// Default Default
const Default = "faas"

func (c *client) CreateGitToken(ctx context.Context, host, token string) error {
	secret := c.client.CoreV1().Secrets(c.k8sNamespace)
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
			Namespace: c.k8sNamespace,
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

// CreateGitSSH CreateGitSSH
func (c *client) CreateGitSSH(ctx context.Context, host, ssh string) error {
	secret := c.client.CoreV1().Secrets(c.k8sNamespace)
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	if tenantID == "" || tenantID == UnexpectedType {
		tenantID = Default
	}
	data := make(map[string][]byte)
	data["known_hosts"] = []byte(host)
	data["ssh-privatekey"] = []byte(ssh)
	s := &v1.Secret{
		Type: v1.SecretTypeSSHAuth,
		ObjectMeta: ctrl.ObjectMeta{
			Name:      tenantID + "-git",
			Namespace: c.k8sNamespace,
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
	sc := c.client.CoreV1().Secrets(c.k8sNamespace)
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
			Namespace: c.k8sNamespace,
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

// Function Function
type Function struct {
	ID        string
	Version   string
	Project   string
	GroupName string
	Git       *Git
	Docker    *Docker
	Builder   string
	ENV       map[string]string
}

// Docker Docker
type Docker struct {
	Host      string
	NameSpace string
	Name      string
}

// Git Git
type Git struct {
	Name string
	Host string
}

const (
	// RESOURCRREF RESOURCRREF
	RESOURCRREF = "kubernetes.pod_name"
	// STEP  STEP
	STEP = "kubernetes.container_name"
	// BuildID BuildID
	BuildID = "quanxiang.faas.build/id"
	// GROUP GROUP
	GROUP = "quanxiang.faas/group"
	// ProjectTAG ProjectTAG
	ProjectTAG = "quanxiang.faas.project/tag"
	// PROJECT PROJECT
	PROJECT = "quanxiang.faas/project"
	// TenentID TenentID
	TenentID = "quanxiang.faas/tenantID"
	// ModuleNAME ModuleNAME
	ModuleNAME = "quanxiang.faas.module/name"
	// BUILD BUILD
	BUILD = "build"
)

func (c *client) Build(ctx context.Context, data *Function) error {
	fn := c.ofn.CoreV1beta1().Functions(c.k8sNamespace)
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
					Url: data.Git.Host + data.GroupName + "/" + data.Project + ".git",
					Credentials: &v1.LocalObjectReference{
						Name: data.Git.Name,
					},
				},
			},
		},
	}
	_, err := fn.Create(ctx, function, metav1.CreateOptions{})
	return err
}

// GetBuilder GetBuilder
func GetBuilder(language string) string {
	return "openfunction/builder-go:latest"
}

// DelFunction DelFunction
type DelFunction struct {
	Name string
}

func (c *client) DelFunction(ctx context.Context, data *DelFunction) error {
	fn := c.ofn.CoreV1beta1().Functions(c.k8sNamespace)
	return fn.Delete(ctx, data.Name, metav1.DeleteOptions{})
}
