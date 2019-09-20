package workshop

import (
	"context"

	"github.com/go-logr/logr"
	workshopv1alpha1 "github.com/kameshsampath/workshop-operator/pkg/apis/kameshs/v1alpha1"
	"github.com/kameshsampath/workshop-operator/pkg/create"
	osappsv1 "github.com/openshift/api/apps/v1"
	isv1 "github.com/openshift/api/image/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileWorkshop) installNexus(reqLogger logr.Logger, workshop *workshopv1alpha1.Workshop) error {
	//Install Nexus

	nexusIS := create.NexusImageStream()
	is := &isv1.ImageStream{}

	err := r.client.Get(context.TODO(),
		types.NamespacedName{Name: create.NexusDeploymentName, Namespace: create.RHDWorkshopNamespace}, is)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating Nexus Image Stream")
		err = r.client.Create(context.TODO(), nexusIS)
		if err != nil {
			return err
		}
		controllerutil.SetControllerReference(workshop, nexusIS, r.scheme)
	}

	nexus := create.Nexus()
	ndc := &osappsv1.DeploymentConfig{}

	err = r.client.Get(context.TODO(),
		types.NamespacedName{Name: create.NexusDeploymentName, Namespace: create.RHDWorkshopNamespace}, ndc)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Deploying Nexus")
		err = r.client.Create(context.TODO(), nexus)
		if err != nil {
			return err
		}
		controllerutil.SetControllerReference(workshop, nexus, r.scheme)
	}

	//Nexus Service
	nexusSvc := create.NexusService()
	err = r.client.Get(context.TODO(),
		types.NamespacedName{Name: create.NexusDeploymentName, Namespace: create.RHDWorkshopNamespace}, nexusSvc)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Deploying Nexus Service")
		err = r.client.Create(context.TODO(), nexusSvc)
		if err != nil {
			return err
		}
		controllerutil.SetControllerReference(workshop, nexusSvc, r.scheme)
	}

	return nil
}
