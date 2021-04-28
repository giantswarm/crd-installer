package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var (
	oneVersion = apiextensionsv1.CustomResourceDefinition{
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
	twoVersions = apiextensionsv1.CustomResourceDefinition{
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
	}
)

func Test_ensureCRDs(t *testing.T) {
	fakeBuilder := fake.NewClientBuilder()
	scheme, err := createScheme()
	require.Nil(t, err)
	fakeBuilder.WithScheme(scheme)

	fakeBuilder.WithRuntimeObjects(oneVersion.DeepCopy())
	fakeClient := fakeBuilder.Build()

	err = ensureCRDs(fakeClient, []*apiextensionsv1.CustomResourceDefinition{twoVersions.DeepCopy()})
	require.Nil(t, err)

	var updatedCRD apiextensionsv1.CustomResourceDefinition
	err = fakeClient.Get(context.Background(), client.ObjectKey{Name: oneVersion.Name}, &updatedCRD)
	require.Nil(t, err)
	require.Len(t, updatedCRD.Spec.Versions, 2)
}
