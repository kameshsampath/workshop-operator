package util

import (
	clientv1alpha1 "github.com/operator-framework/operator-lifecycle-manager/pkg/api/client/clientset/versioned/typed/operators/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("controller_workshop")

//GetOlmClient gets the OLM Client
func GetOlmClient() (*clientv1alpha1.OperatorsV1alpha1Client, error) {
	config, err := config.GetConfig()

	if err != nil {
		log.Error(err, "Error getting configuration")
		return nil, err
	}

	clientSet, err := clientv1alpha1.NewForConfig(config)

	if err != nil {
		log.Error(err, "Error getting client set")
		return nil, err
	}

	return clientSet, nil
}
