package main

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"flag"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"reflect"

	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	apiyaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	crdDirectory = flag.String("dir", "", "")
)

func createClient() (client.Client, error) {
	restConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientScheme := runtime.NewScheme()
	err = v1.AddToScheme(clientScheme)
	if err != nil {
		return nil, err
	}

	k8sClient, err := client.New(restConfig, client.Options{
		Scheme: clientScheme,
	})
	if err != nil {
		return nil, err
	}

	return k8sClient, nil
}

func decodeCRDs(readCloser io.ReadCloser) ([]*v1.CustomResourceDefinition, error) {
	reader := apiyaml.NewYAMLReader(bufio.NewReader(readCloser))
	decoder := scheme.Codecs.UniversalDeserializer()

	defer func(contentReader io.ReadCloser) {
		err := readCloser.Close()
		if err != nil {
			panic(err)
		}
	}(readCloser)

	crdGVK := schema.GroupVersionKind{
		Group:   "apiextensions.k8s.io",
		Version: "v1",
		Kind:    "CustomResourceDefinition",
	}
	var crds []*v1.CustomResourceDefinition
	for {
		doc, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, err
		}

		//  Skip over empty documents, i.e. a leading `---`
		if len(bytes.TrimSpace(doc)) == 0 {
			continue
		}

		var crd v1.CustomResourceDefinition
		_, decodedGVK, err := decoder.Decode(doc, nil, &crd)
		if err != nil {
			return nil, err
		} else if *decodedGVK != crdGVK {
			continue
		}
		crds = append(crds, &crd)
	}

	return crds, nil
}

func findCRDs(dir string) ([]*v1.CustomResourceDefinition, error) {
	var crds []*v1.CustomResourceDefinition
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		documentReader, err := os.Open(path)
		if err != nil {
			return err
		}

		documentCRDs, err := decodeCRDs(documentReader)
		if err != nil {
			return err
		}

		crds = append(crds, documentCRDs...)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return crds, nil
}

func run() error {
	crds, err := findCRDs(*crdDirectory)
	if err != nil {
		return err
	}

	k8sClient, err := createClient()
	if err != nil {
		return err
	}

	err = ensureCRDs(k8sClient, crds)
	if err != nil {
		return err
	}

	return nil
}

func ensureCRDs(k8sClient client.Client, crds []*v1.CustomResourceDefinition) error {
	ctx := context.Background()
	for _, crd := range crds {
		var existing v1.CustomResourceDefinition
		err := k8sClient.Get(ctx, client.ObjectKey{Name: crd.Name}, &existing)
		if apierrors.IsNotFound(err) {
			err := k8sClient.Create(ctx, crd)
			if err != nil {
				return err
			}

			continue
		} else if err != nil {
			return err
		}

		if reflect.DeepEqual(crd.Spec, existing.Spec) {
			continue
		}

		existing.Spec = crd.Spec
		err = k8sClient.Update(ctx, &existing)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	flag.Parse()
	if *crdDirectory == "" {
		log.Fatal("-dir must be defined")
	}

	if err := run(); errors.Is(err, fs.ErrNotExist) {
		log.Fatal("dir not found")
	} else if err != nil {
		log.Fatal(err)
	}
}
