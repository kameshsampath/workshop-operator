package workshop

import (
	"context"

	"github.com/go-logr/logr"
	workshopv1alpha1 "github.com/kameshsampath/workshop-operator/pkg/apis/kameshs/v1alpha1"
	"github.com/kameshsampath/workshop-operator/pkg/create"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileWorkshop) createProjects(reqLogger logr.Logger, workshop *workshopv1alpha1.Workshop) error {
	projects := create.WorkshopProjects(workshop.Spec)
	for _, p := range projects {
		err := r.client.Get(context.TODO(),
			types.NamespacedName{Name: p.Name}, p)
		if err != nil {
			if errors.IsNotFound(err) {
				reqLogger.Info("Creating Workshop Project", "Project Name", p.Name)
				err := r.client.Create(context.TODO(), p)
				if err != nil {
					return err
				}
				controllerutil.SetControllerReference(workshop, p, r.scheme)
			} else {
				reqLogger.Info("Creating Workshop Project", "Project Name", p.Name, "Already Exists")
			}
		}
	}
	return nil
}
