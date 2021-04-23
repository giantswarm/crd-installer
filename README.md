[![CircleCI](https://circleci.com/gh/giantswarm/crd-installer.svg?style=shield)](https://circleci.com/gh/giantswarm/crd-installer)

# crd-installer

This program is intended to be used as an init-container in a Kubernetes pod to ensure CRDs have been installed or 
updated before the main container runs. It accepts a single flag, `-dir` which should point to a directory containing
one or more YAML-formatted CRDs to be installed. It uses in-cluster credentials to authenticate to the Kubernetes API
using a service account token.
