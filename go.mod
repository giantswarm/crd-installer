module github.com/giantswarm/crd-installer

go 1.16

require (
	github.com/stretchr/testify v1.6.1
	k8s.io/apiextensions-apiserver v0.20.6
	k8s.io/apimachinery v0.20.6
	k8s.io/client-go v0.20.6
	sigs.k8s.io/controller-runtime v0.8.3
)

// v3.3.13 is required by bketelsen/crypt. Can remove this replace when updated.
replace github.com/coreos/etcd v3.3.13+incompatible => github.com/coreos/etcd v3.3.25+incompatible
