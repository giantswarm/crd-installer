package main

import (
	"context"
	"io"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/stretchr/testify/require"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
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
	crds, err := decodeCRDs(readCloser)
	require.Nil(t, err)
	require.Len(t, crds, 1)
}

func Test_findCRDs(t *testing.T) {
	filesystem := fstest.MapFS{
		"a.yaml": &fstest.MapFile{
			Data:    []byte(testCRDDocument),
			Mode:    0700,
			ModTime: time.Now(),
			Sys:     nil,
		},
	}
	crds, err := findCRDs(filesystem, ".")
	require.Nil(t, err)
	require.Len(t, crds, 1)
}

func Test_ensureCRDs(t *testing.T) {
	fakeBuilder := fake.NewClientBuilder()
	scheme := runtime.NewScheme()
	err := apiextensionsv1.AddToScheme(scheme)
	require.Nil(t, err)
	fakeBuilder.WithScheme(scheme)
	existingCRD := apiextensionsv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: "existing",
		},
		Spec: apiextensionsv1.CustomResourceDefinitionSpec{
			Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
				{
					Name: "v1alpha1",
				},
			},
		},
	}
	fakeBuilder.WithRuntimeObjects(&existingCRD)
	fakeClient := fakeBuilder.Build()

	err = ensureCRDs(fakeClient, []*apiextensionsv1.CustomResourceDefinition{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "existing",
			},
			Spec: apiextensionsv1.CustomResourceDefinitionSpec{
				Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
					{
						Name: "v1alpha1",
					},
					{
						Name: "v1beta1",
					},
				},
			},
		},
	})
	require.Nil(t, err)

	var updatedCRD apiextensionsv1.CustomResourceDefinition
	err = fakeClient.Get(context.Background(), client.ObjectKey{Name: "existing"}, &updatedCRD)
	require.Nil(t, err)
	require.Len(t, updatedCRD.Spec.Versions, 2)
}
