package update

import (
	"encoding/json"

	"github.com/kameshsampath/workshop-operator/pkg/create"
	"github.com/kameshsampath/workshop-operator/pkg/util"

	"github.com/go-logr/logr"
	workshopv1alpha1 "github.com/kameshsampath/workshop-operator/pkg/apis/kameshs/v1alpha1"
	types "k8s.io/apimachinery/pkg/types"
)

//CheForTokenExchange enables che Operator install Che Keycloak with Token Exchange
func CheForTokenExchange(spec workshopv1alpha1.WorkshopSpec, log logr.Logger) error {
	cheOpImg := spec.Stack.Che.OperatorImage
	clientSet, err := util.GetOlmClient()

	if err != nil {
		return err
	}

	patches := []util.PatchAsString{
		util.PatchAsString{
			Op:    "replace",
			Path:  "/metadata/annotations/containerImage",
			Value: cheOpImg,
		},
		util.PatchAsString{
			Op:    "replace",
			Path:  "/spec/install/spec/deployments/0/spec/template/spec/containers/0/image",
			Value: cheOpImg,
		},
	}

	payloadBytes, _ := json.Marshal(patches)

	result, err := clientSet.
		ClusterServiceVersions(create.CheInstallNamespace).
		Patch("eclipse-che.v"+spec.Stack.Che.CheVersion, types.JSONPatchType, payloadBytes)

	if err != nil {
		return err
	}

	if result != nil {
		log.Info("Successfully patched Eclipse Che for Token Exchange")
	}

	return nil
}
