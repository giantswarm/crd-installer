package main

import (
	"context"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func createScheme() (*runtime.Scheme, error) {
	scheme := runtime.NewScheme()
	err := v1.AddToScheme(scheme)
	if err != nil {
		return nil, errors.Wrap(err, "add crd v1 to scheme error")
	}
	return scheme, nil
}

func createClient(clientScheme *runtime.Scheme, kubeconfigPath string) (client.Client, error) {
	var restConfig *rest.Config
	if kubeconfigPath == "" {
		globalLogger.Info("kubeconfig flag not specified, using in-cluster credentials")
		var err error
		restConfig, err = rest.InClusterConfig()
		if err != nil {
			return nil, errors.Wrap(err, "in-cluster rest config error")
		}
	} else {
		globalLogger.Info("reading credentials from kubeconfig", zap.String("kubeconfig", kubeconfigPath))
		var err error
		restConfig, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
			&clientcmd.ConfigOverrides{}).ClientConfig()
		if err != nil {
			return nil, errors.Wrap(err, "rest config from kubeconfig error")
		}
	}

	k8sClient, err := client.New(restConfig, client.Options{
		Scheme: clientScheme,
	})
	if err != nil {
		return nil, errors.Wrap(err, "create client error")
	}

	return k8sClient, nil
}

func ensureCRDs(k8sClient client.Client, crds []*v1.CustomResourceDefinition) error {
	globalLogger.Info("ensuring all crds exist and are and up-to-date")

	ctx := context.Background()
	for _, crd := range crds {
		logger := globalLogger.With(zap.String("crd", crd.Name))
		logger.Info("ensuring")

		var existing v1.CustomResourceDefinition
		err := k8sClient.Get(ctx, client.ObjectKey{Name: crd.Name}, &existing)
		if apierrors.IsNotFound(err) {
			err := k8sClient.Create(ctx, crd)
			if err != nil {
				return errors.Wrap(err, "create error")
			}

			logger.Info("created")

			continue
		} else if err != nil {
			return errors.Wrap(err, "get error")
		}

		diff := cmp.Diff(crd.Spec, existing.Spec)
		if diff == "" {
			logger.Info("spec matches")
			continue
		}

		logger.Info("spec differs", zap.String("diff", diff))

		existing.Spec = crd.Spec
		err = k8sClient.Update(ctx, &existing)
		if err != nil {
			return errors.Wrap(err, "update error")
		}

		logger.Info("updated")

		err = k8sClient.Get(ctx, client.ObjectKey{Name: crd.Name}, &existing)
		if err != nil {
			return errors.Wrap(err, "check status error")
		}

		for _, condition := range existing.Status.Conditions {
			if condition.Status == v1.ConditionFalse {
				logger.Error("found error condition in status",
					zap.String("condition", string(condition.Type)),
					zap.String("status", string(condition.Status)))
				return errors.New("error condition in status when ensuring crd")
			}
		}
	}

	return nil
}
