package k8s

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/openfunction/apis/core/v1beta1"
	"github.com/openfunction/pkg/client/clientset/versioned"
	ginheader "github.com/quanxiang-cloud/cabin/tailormade/header"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/utils/pointer"

	pipeline "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	tektonClient "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	ksvc "knative.dev/serving/pkg/apis/serving/v1"
	serving "knative.dev/serving/pkg/client/clientset/versioned/typed/serving/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

// Client Client
type Client interface {
	CreateGitToken(ctx context.Context, host, token string) error
	CreateGitSSH(ctx context.Context, host, ssh string) error
	CreateDocker(ctx context.Context, host, username, secret string) error
	Build(ctx context.Context, data *Function) error
	DelFunction(ctx context.Context, data *DelFunction) error
	CreateServing(ctx context.Context, fn *Function) error
	DelServing(ctx context.Context, fn *Function) error
	RegistAPI(ctx context.Context, fn *Function, appId string) error
	DeleteReigstRun(ctx context.Context, name string) error
}

type client struct {
	client          *kubernetes.Clientset
	tekton          *tektonClient.Clientset
	ofn             versioned.Interface
	serving         serving.ServiceInterface
	k8sNamespace    string
	dockerNamespace string
}

// NewClient NewClient
func NewClient(namespace string) Client {
	config := ctrl.GetConfigOrDie()
	clientset := kubernetes.NewForConfigOrDie(config)
	ofn := versioned.NewForConfigOrDie(config)
	return &client{
		client:       clientset,
		k8sNamespace: namespace,
		ofn:          ofn,
		serving:      serving.NewForConfigOrDie(config).Services(namespace),
		tekton:       tektonClient.NewForConfigOrDie(config),
	}
}

const (
	// UnexpectedType UnexpectedType
	UnexpectedType = "unexpected type"

	// Default Default
	Default = "faas"

	// GITTEKTON GITTEKTON
	GITTEKTON = "tekton.dev/git-0"
)

func (c *client) CreateGitToken(ctx context.Context, host, token string) error {
	secret := c.client.CoreV1().Secrets(c.k8sNamespace)
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	if tenantID == "" || tenantID == UnexpectedType {
		tenantID = Default
	}
	data := make(map[string][]byte)
	data["host"] = []byte(host)
	data["token"] = []byte(token)

	tekton := make(map[string]string)
	tekton[GITTEKTON] = host
	s := &v1.Secret{
		Type: v1.SecretTypeOpaque,
		ObjectMeta: ctrl.ObjectMeta{
			Name:        tenantID + "-git",
			Namespace:   c.k8sNamespace,
			Annotations: tekton,
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

	tekton := make(map[string]string)
	tekton[GITTEKTON] = host

	s := &v1.Secret{
		Type: v1.SecretTypeSSHAuth,
		ObjectMeta: ctrl.ObjectMeta{
			Name:        tenantID + "-git",
			Namespace:   c.k8sNamespace,
			Annotations: tekton,
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

func (c *client) CreateServing(ctx context.Context, fn *Function) error {
	ksvc := &ksvc.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      GenName(fn, true),
			Namespace: c.k8sNamespace,
		},
		Spec: ksvc.ServiceSpec{
			ConfigurationSpec: ksvc.ConfigurationSpec{
				Template: ksvc.RevisionTemplateSpec{
					Spec: ksvc.RevisionSpec{
						PodSpec: v1.PodSpec{
							ImagePullSecrets: []v1.LocalObjectReference{
								{Name: fn.Docker.Name},
							},
							Containers: []v1.Container{{
								Name: "serving",
								// Image: fn.Docker.Host + fn.Docker.NameSpace + strings.ToLower(fn.GroupName) + "-" + fn.Project + ":" + fn.Version,
								Image: fn.Docker.Host + fn.Docker.NameSpace + strings.ToLower(fn.GroupName) + "-" + fn.Project + ":" + fn.Version,
								Ports: []v1.ContainerPort{{
									ContainerPort: 8080,
								}},
								Env: genEnv(c, fn),
							}},
						},
					},
				},
			},
		},
	}

	_, err := c.serving.Create(ctx, ksvc, metav1.CreateOptions{})
	return err
}

func (c *client) DelServing(ctx context.Context, fn *Function) error {
	return c.serving.Delete(ctx,
		GenName(fn, true),
		metav1.DeleteOptions{})
}

func genEnv(c *client, fn *Function) []v1.EnvVar {
	env := make([]v1.EnvVar, 0, len(fn.ENV))
	for k, v := range fn.ENV {
		env = append(env, v1.EnvVar{
			Name:  k,
			Value: v,
		})
	}
	env = append(env,
		v1.EnvVar{
			Name:  "FUNC_CONTEXT",
			Value: fmt.Sprintf("{\"name\":\"%s\",\"version\":\"v2.0.0\",\"runtime\":\"Knative\",\"port\":\"8080\"}", GenName(fn, true)),
		},
		v1.EnvVar{
			Name: "POD_NAME",
			ValueFrom: &v1.EnvVarSource{
				FieldRef: &v1.ObjectFieldSelector{
					APIVersion: "v1",
					FieldPath:  "metadata.name",
				},
			},
		},
		v1.EnvVar{
			Name:  "POD_NAMESPACE",
			Value: c.k8sNamespace,
		})

	return env
}

func GenName(fn *Function, reverse bool, prefix ...string) (ret string) {
	template := "%s-%s-%s"
	if len(prefix) != 0 {
		template = prefix[0] + template
	}

	if !reverse {
		ret = fmt.Sprintf(template, fn.GroupName, fn.Project, fn.Version)
	} else {
		ret = fmt.Sprintf(template, fn.Version, fn.Project, fn.GroupName)
	}
	ret = strings.ToLower(ret)
	return
}

func ReverseName(name string) (string, error) {
	first := strings.Index(name, "-")
	if first == -1 {
		return name, fmt.Errorf("invalid name")
	}
	last := strings.LastIndex(name, "-")
	if first != last {
		return fmt.Sprintf("%s-%s-%s", name[last+1:], name[first+1:last], name[:first]), nil
	}
	return fmt.Sprintf("%s-%s", name[last+1:], name[:first]), nil
}

// TODO: check host
func genGitRepo(fn *Function) string {
	host, group, project := fn.Git.Host, fn.GroupName, fn.Project
	if index := strings.LastIndex(host, "/"); index == len(host)-1 {
		host = host[:index]
	}
	return fmt.Sprintf("%s/%s/%s.git", host, group, project)
}

func (c *client) RegistAPI(ctx context.Context, fn *Function, appId string) error {
	pipeRun := &pipeline.PipelineRun{
		ObjectMeta: metav1.ObjectMeta{
			Name: GenName(fn, true),
		},
		Spec: pipeline.PipelineRunSpec{
			PipelineRef: &pipeline.PipelineRef{
				Name: "register-polyapi",
			},
			Params: []pipeline.Param{
				{
					Name:  "SOURCE_URL",
					Value: *pipeline.NewArrayOrString(genGitRepo(fn)),
				},
				{
					Name:  "PROJECT_NAME",
					Value: *pipeline.NewArrayOrString(fn.Project),
				},
				{
					Name:  "OPERATE_ID",
					Value: *pipeline.NewArrayOrString(GenName(fn, false)),
				},
				{
					Name: "HOST",
					// TODO: hard code
					Value: *pipeline.NewArrayOrString("localhost:9999"),
				},
				{
					Name:  "APPID",
					Value: *pipeline.NewArrayOrString(appId),
				},
			},
			ServiceAccountName: "builder",
			Workspaces: []pipeline.WorkspaceBinding{
				{
					Name:     "source-ws",
					EmptyDir: &v1.EmptyDirVolumeSource{},
				},
			},
		},
	}
	_, err := c.tekton.TektonV1beta1().PipelineRuns(c.k8sNamespace).Create(ctx, pipeRun, metav1.CreateOptions{})
	return err
}

func (c *client) DeleteReigstRun(ctx context.Context, name string) error {
	return c.tekton.TektonV1beta1().PipelineRuns(c.k8sNamespace).Delete(ctx, name, metav1.DeleteOptions{})
}
