module github.com/giantswarm/crd-installer

go 1.16

require (
	github.com/google/go-cmp v0.5.9
	github.com/imdario/mergo v0.3.10 // indirect
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.8.2
	go.uber.org/zap v1.24.0
	k8s.io/apiextensions-apiserver v0.26.1
	k8s.io/apimachinery v0.26.1
	k8s.io/client-go v0.26.1
	sigs.k8s.io/controller-runtime v0.14.5
)

// v3.3.13 is required by bketelsen/crypt. Can remove this replace when updated.
replace (
	github.com/coreos/etcd v3.3.13+incompatible => github.com/coreos/etcd v3.3.25+incompatible
	github.com/dgrijalva/jwt-go => github.com/dgrijalva/jwt-go/v4 v4.0.0-preview1
)
