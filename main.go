package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	crdDirectory   = flag.String("dir", "", "")
	kubeconfigPath = flag.String("kubeconfig", "", "")

	globalLogger, _ = zap.NewDevelopment(zap.AddStacktrace(zap.FatalLevel))
)

func run() error {
	var crds []*v1.CustomResourceDefinition
	{
		globalLogger.Info("reading CRDs", zap.String("directory", *crdDirectory))

		var err error
		if crds, err = findCRDs(*crdDirectory); err != nil {
			return err
		} else if len(crds) == 0 {
			globalLogger.Info("found no CRDs, exiting")
			return nil
		}

		var crdNames []string
		for _, crd := range crds {
			crdNames = append(crdNames, crd.Name)
		}
		globalLogger.Info("found CRDs",
			zap.Int("count", len(crds)),
			zap.Strings("names", crdNames))
	}

	var k8sClient client.Client
	{
		scheme, err := createScheme()
		if err != nil {
			return err
		}

		if k8sClient, err = createClient(scheme, *kubeconfigPath); err != nil {
			return err
		}
	}

	if err := ensureCRDs(k8sClient, crds); err != nil {
		return err
	}

	globalLogger.Info("completed successfully")
	return nil
}

func main() {
	flag.Parse()
	if *crdDirectory == "" {
		globalLogger.Error("-dir must be defined")
		os.Exit(1)
	}

	err := run()
	if errors.Is(err, fs.ErrNotExist) {
		globalLogger.Error("dir not found")
		os.Exit(1)
	} else if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}
}
