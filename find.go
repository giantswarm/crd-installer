package main

import (
	"bufio"
	"bytes"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	apiyaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
)

var (
	crdGroupVersionKind = schema.GroupVersionKind{
		Group:   "apiextensions.k8s.io",
		Version: "v1",
		Kind:    "CustomResourceDefinition",
	}
)

func decodeCRDs(filename string, readCloser io.ReadCloser) ([]*v1.CustomResourceDefinition, error) {
	logger := globalLogger.With(zap.String("filename", filename))
	reader := apiyaml.NewYAMLReader(bufio.NewReader(readCloser))
	decoder := scheme.Codecs.UniversalDeserializer()

	defer func(contentReader io.ReadCloser) {
		err := readCloser.Close()
		if err != nil {
			panic(errors.Wrap(err, "deferred close error"))
		}
	}(readCloser)

	var crds []*v1.CustomResourceDefinition
	for {
		doc, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, errors.Wrap(err, "read error")
		}

		//  Skip over empty documents, i.e. a leading `---`
		if len(bytes.TrimSpace(doc)) == 0 {
			continue
		}

		var crd v1.CustomResourceDefinition
		_, decodedGVK, err := decoder.Decode(doc, nil, &crd)
		if err != nil {
			return nil, errors.Wrap(err, "decode error")
		} else if *decodedGVK != crdGroupVersionKind {
			logger.Info("skipping non-CRD",
				zap.String("group", decodedGVK.Group),
				zap.String("version", decodedGVK.Version),
				zap.String("kind", decodedGVK.Kind))
			continue
		}
		logger.Info("found CRD", zap.String("name", crd.Name))
		crds = append(crds, &crd)
	}

	return crds, nil
}

func findCRDs(dir string) ([]*v1.CustomResourceDefinition, error) {
	var crds []*v1.CustomResourceDefinition
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return errors.Wrap(err, "walk dir func error")
		}
		if d.IsDir() {
			globalLogger.Info("skipping directory", zap.String("path", path))
			return nil
		} else if filepath.Ext(path) != ".yaml" {
			globalLogger.Info("skipping non-YAML file", zap.String("path", path))
			return nil
		}

		documentReader, err := os.Open(path)
		if err != nil {
			return errors.Wrap(err, "open error")
		}

		documentCRDs, err := decodeCRDs(path, documentReader)
		if err != nil {
			return err
		}

		crds = append(crds, documentCRDs...)
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "walk dir error")
	}

	return crds, nil
}
