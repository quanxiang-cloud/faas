module github.com/quanxiang-cloud/faas

go 1.16

require (
	github.com/gin-gonic/gin v1.7.7
	github.com/go-logr/logr v1.2.2
	github.com/go-logr/zapr v1.2.2
	github.com/openfunction v0.0.0-00010101000000-000000000000
	github.com/quanxiang-cloud/cabin v0.0.6
	go.uber.org/zap v1.19.1
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/gorm v1.22.4
	k8s.io/api v0.23.5
	k8s.io/apimachinery v0.23.5
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	k8s.io/utils v0.0.0-20211116205334-6203023598ed
	sigs.k8s.io/controller-runtime v0.11.2
)

replace (
	github.com/go-logr/logr => github.com/go-logr/logr v1.2.0
	github.com/openfunction => github.com/OpenFunction/OpenFunction v0.6.0
	//github.com/quanxiang-cloud/organizations => ../../../github.com/quanxiang-cloud/organizations
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.21.4
	k8s.io/client-go => k8s.io/client-go v0.23.5
	k8s.io/klog/v2 => k8s.io/klog/v2 v2.30.0
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.11.2
)

require (
	github.com/go-redis/redis/v8 v8.11.4
	github.com/jinzhu/now v1.1.4 // indirect
	github.com/olivere/elastic/v7 v7.0.30
	github.com/quanxiang-cloud/organizations v1.0.3
	//github.com/quanxiang-cloud/organizations v0.0.0-00010101000000-000000000000
	// v1.0.3
	github.com/xanzy/go-gitlab v0.63.0
)
