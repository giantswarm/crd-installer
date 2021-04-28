package main

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const testCRDDocument = `apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: crontabs.stable.example.com
spec:
  group: stable.example.com
  versions:
    - name: v1
      served: true
      storage: true
      schema:
        openAPIV3Schema:
          type: object
          properties:
            spec:
              type: object
              properties:
                cronSpec:
                  type: string
                image:
                  type: string
                replicas:
                  type: integer
  scope: Namespaced
  names:
    plural: crontabs
    singular: crontab
    kind: CronTab
    shortNames:
    - ct`

func Test_decodeCRDs(t *testing.T) {
	readCloser := io.NopCloser(strings.NewReader(testCRDDocument))
	crds, err := decodeCRDs("test", readCloser)
	require.Nil(t, err)
	require.Len(t, crds, 1)
}
