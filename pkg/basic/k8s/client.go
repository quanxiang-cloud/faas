package k8s

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/openfunction/apis/core/v1beta1"
	"github.com/openfunction/pkg/client/clientset/versioned"
	ginheader "github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/faas/pkg/basic/strutil"
	"github.com/quanxiang-cloud/faas/pkg/config"
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
	GetBuilder(language string, version string) (string, error)
}

type client struct {
	client       *kubernetes.Clientset
	tekton       *tektonClient.Clientset
	ofn          versioned.Interface
	serving      serving.ServiceInterface
	k8sNamespace string
	buildImages  map[string]string
}

// NewClient NewClient
func NewClient(config *config.Config, namespace string) Client {
	ctrlConfig := ctrl.GetConfigOrDie()
	clientset := kubernetes.NewForConfigOrDie(ctrlConfig)
	ofn := versioned.NewForConfigOrDie(ctrlConfig)

	return &client{
		client:       clientset,
		k8sNamespace: namespace,
		ofn:          ofn,
		serving:      serving.NewForConfigOrDie(ctrlConfig).Services(namespace),
		tekton:       tektonClient.NewForConfigOrDie(ctrlConfig),
		buildImages:  config.BuildImages,
	}
}

func (c *client) CreateGitToken(ctx context.Context, host, token string) error {
	secret := c.client.CoreV1().Secrets(c.k8sNamespace)
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	if tenantID == "" || tenantID == TenantUnexpectedType {
		tenantID = TenantDefault
	}
	data := make(map[string][]byte)
	data["host"] = []byte(host)
	data["token"] = []byte(token)

	tekton := make(map[string]string)
	tekton[GitTekton] = host
	s := &v1.Secret{
		Type: v1.SecretTypeOpaque,
		ObjectMeta: ctrl.ObjectMeta{
			Name:        strutil.GenName(tenantID, SecretGitSuffix),
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
	if tenantID == "" || tenantID == TenantUnexpectedType {
		tenantID = TenantDefault
	}
	data := make(map[string][]byte)
	data["known_hosts"] = []byte(host)
	data["ssh-privatekey"] = []byte(ssh)

	tekton := make(map[string]string)
	tekton[GitTekton] = host

	s := &v1.Secret{
		Type: v1.SecretTypeSSHAuth,
		ObjectMeta: ctrl.ObjectMeta{
			Name:        strutil.GenName(tenantID, SecretGitSuffix),
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
	if tenantID == "" || tenantID == TenantUnexpectedType {
		tenantID = TenantDefault
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
			Name:      strutil.GenName(tenantID, SecretDockerSuffix),
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

func (c *client) Build(ctx context.Context, data *Function) error {
	fn := c.ofn.CoreV1beta1().Functions(c.k8sNamespace)
	function := &v1beta1.Function{
		ObjectMeta: metav1.ObjectMeta{
			Name: strutil.GenName(data.GroupName, data.Project, data.Version),
		},
		Spec: v1beta1.FunctionSpec{
			Version: &data.Version,
			Image:   strutil.JoinImage(data.Version, data.Docker.Host, data.Docker.NameSpace, strutil.GenName(data.GroupName, data.Project)),
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
					Url: strutil.JoinGIT(data.Git.Host, data.GroupName, data.Project),
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
func (c *client) GetBuilder(language string, version string) (string, error) {
	image, ok := c.buildImages[language+version]
	if !ok {
		return "", fmt.Errorf("the language(%s) is not supported", language+version)
	}
	return image, nil
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
			Name:      strutil.GenName(fn.Version, fn.Project, fn.GroupName),
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
								Name:  KsvcDefaultContainerName,
								Image: strutil.JoinImage(fn.Version, fn.Docker.Host, fn.Docker.NameSpace, strutil.GenName(fn.GroupName, fn.Project)),
								Ports: []v1.ContainerPort{{
									ContainerPort: KsvcDefaultContainerPort,
								}},
								Env: genServingEnv(c, fn),
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
		strutil.GenName(fn.Version, fn.Project, fn.GroupName),
		metav1.DeleteOptions{})
}

func genServingEnv(c *client, fn *Function) []v1.EnvVar {
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
			Value: fmt.Sprintf("{\"name\":\"%s\",\"version\":\"v2.0.0\",\"runtime\":\"Knative\",\"port\":\"8080\", \"prePlugins\":[\"plugin-quanxiang-lowcode-client\"]}", strutil.GenName(fn.Version, fn.Project, fn.GroupName)),
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

func (c *client) RegistAPI(ctx context.Context, fn *Function, appId string) error {
	// TODO:  replace pipelinerun with taskrun
	pipeRun := &pipeline.PipelineRun{
		ObjectMeta: metav1.ObjectMeta{
			Name: strutil.GenName(fn.Version, fn.Project, fn.GroupName),
		},
		Spec: pipeline.PipelineRunSpec{
			PipelineRef: &pipeline.PipelineRef{
				Name: RegisterAPIPipeline,
			},
			Params: []pipeline.Param{
				{
					Name:  RegisterAPIParamSourceURL,
					Value: *pipeline.NewArrayOrString(strutil.JoinGIT(fn.Git.Host, fn.GroupName, fn.Project)),
				},
				{
					Name:  RegisterAPIParamProject,
					Value: *pipeline.NewArrayOrString(fn.Project),
				},
				{
					Name:  RegisterAPIParamPOperate,
					Value: *pipeline.NewArrayOrString(strutil.GenName(fn.GroupName, fn.Project, fn.Version)),
				},
				{
					Name:  RegisterAPIParamAppID,
					Value: *pipeline.NewArrayOrString(appId),
				},
			},
			ServiceAccountName: RegisterAPIServiceAccount,
			Workspaces: []pipeline.WorkspaceBinding{
				{
					Name:     RegisterAPIWorkSpace,
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
