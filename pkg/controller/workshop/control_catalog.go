package workshop

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/kameshsampath/workshop-operator/pkg/create"
	"k8s.io/apimachinery/pkg/api/errors"

	workshopv1alpha1 "github.com/kameshsampath/workshop-operator/pkg/apis/kameshs/v1alpha1"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

//Create the CatalogSource Config
func (r *ReconcileWorkshop) createCatalogSourceConfigs(reqLogger logr.Logger, workshop *workshopv1alpha1.Workshop) error {
	stackCSCs := create.WorkshopOperatorsCatalog(workshop.Spec)
	for _, csc := range stackCSCs {
		err := r.client.Get(context.TODO(),
			types.NamespacedName{Name: csc.Name, Namespace: create.CSCNS}, csc)

		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info(fmt.Sprintf("Creating CatalogSourceConfig : %s", csc.Name))
			err = r.client.Create(context.TODO(), csc)
			if err != nil {
				return err
			}
			controllerutil.SetControllerReference(workshop, csc, r.scheme)
		} else {
			reqLogger.Info(fmt.Sprintf("Updating CatalogSourceConfig : %s", csc.Name))
			err = r.client.Update(context.TODO(), csc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
