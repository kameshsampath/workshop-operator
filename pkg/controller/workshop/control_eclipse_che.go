package workshop

import (
	"context"

	cheorgv1 "github.com/eclipse/che-operator/pkg/apis/org/v1"
	"github.com/go-logr/logr"
	workshopv1alpha1 "github.com/kameshsampath/workshop-operator/pkg/apis/kameshs/v1alpha1"
	"github.com/kameshsampath/workshop-operator/pkg/create"
	"github.com/kameshsampath/workshop-operator/pkg/update"
	olmv1alpha1 "github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileWorkshop) installAndConfigureEclipseChe(reqLogger logr.Logger, workshop *workshopv1alpha1.Workshop) error {
	//Install Eclipse Che
	cheProject := create.CheProject()
	err := r.client.Get(context.TODO(),
		types.NamespacedName{Name: create.CheInstallNamespace}, cheProject)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating Che project")
		err = r.client.Create(context.TODO(), cheProject)
		if err != nil {
			return err
		}
		controllerutil.SetControllerReference(workshop, cheProject, r.scheme)
	}

	cheCSC := create.CheCatalogSourceConfig()
	err = r.client.Get(context.TODO(),
		types.NamespacedName{Name: create.CheCSCName, Namespace: create.CSCNS}, cheCSC)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating Che CatalogSource Config")
		err = r.client.Create(context.TODO(), cheCSC)
		if err != nil {
			return err
		}
		controllerutil.SetControllerReference(workshop, cheCSC, r.scheme)
	}

	cheOG := create.CheOperatorGroup()
	err = r.client.Get(context.TODO(),
		types.NamespacedName{Name: create.CheOperatorGroupName, Namespace: create.CheInstallNamespace}, cheOG)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating Che Operator Group ")
		err = r.client.Create(context.TODO(), cheOG)
		if err != nil {
			return err
		}
		controllerutil.SetControllerReference(workshop, cheOG, r.scheme)
	}

	cheSub := create.CheSubscription(workshop.Spec)
	err = r.client.Get(context.TODO(),
		types.NamespacedName{Name: create.CheSubscriptionName, Namespace: create.CheInstallNamespace}, cheSub)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating Che Subscription")
		err = r.client.Create(context.TODO(), cheSub)
		if err != nil {
			return err
		}
		controllerutil.SetControllerReference(workshop, cheSub, r.scheme)
	}

	//TEMPORARY: Patch the pod with Keycloak enabled advanced features enabled image

	if workshop.Spec.Stack.Che.EnableTokenExchange {
		log.Info("Enabling Token Exchange")
		cheCsv := &olmv1alpha1.ClusterServiceVersion{}

		err := r.client.Get(context.TODO(),
			types.NamespacedName{Name: "eclipse-che.v" + workshop.Spec.Stack.Che.CheVersion, Namespace: create.CheInstallNamespace}, cheCsv)

		if err != nil {
			if errors.IsNotFound(err) {
				log.Error(err, "Error loading Eclipse Che CSV")
			}
			return err
		}

		err = update.CheForTokenExchange(workshop.Spec, reqLogger)
		if err != nil {
			return err
		}
	}

	che := create.CheCluster(workshop.Spec)
	eche := &cheorgv1.CheCluster{}

	err = r.client.Get(context.TODO(),
		types.NamespacedName{Name: create.CheSubscriptionName, Namespace: create.CheInstallNamespace}, eche)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating Che Cluster")
		err = r.client.Create(context.TODO(), che)
		if err != nil {
			return err
		}
		controllerutil.SetControllerReference(workshop, che, r.scheme)
	}

	return nil
}
