package k8s

// Function Function
type Function struct {
	ID           string
	Version      string
	Project      string
	ProjectTitle string
	GroupName    string
	Git          *Git
	Docker       *Docker
	Builder      string
	ENV          map[string]string
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

// const of faas
const (
	ResourceRef = "kubernetes.pod_name"
	Step        = "kubernetes.container_name"
	GitTekton   = "tekton.dev/git-0"

	TenantUnexpectedType = "unexpected type"
	TenantDefault        = "faas"

	SecretGitSuffix    = "git"
	SecretDockerSuffix = "docker"

	KsvcDefaultContainerName = "serving"
	KsvcDefaultContainerPort = 8080

	RegisterAPIPipeline = "register-polyapi"

	RegisterAPIParamSourceURL    = "SOURCE_URL"
	RegisterAPIParamProject      = "PROJECT_NAME"
	RegisterAPIParamProjectTitle = "PROJECT_TITLE"
	RegisterAPIParamGroup        = "GROUP_NAME"
	RegisterAPIParamPOperate     = "OPERATE_ID"
	RegisterAPIParamAppID        = "APPID"

	RegisterAPIServiceAccount = "builder"
	RegisterAPIWorkSpace      = "source-ws"
)
