package workshop

import (
	"context"
	"fmt"

	"github.com/kameshsampath/workshop-operator/pkg/create"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	workshopv1alpha1 "github.com/kameshsampath/workshop-operator/pkg/apis/kameshs/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	oauthv1 "github.com/openshift/api/config/v1"
	projectv1 "github.com/openshift/api/project/v1"
	userv1 "github.com/openshift/api/user/v1"
	rbac "k8s.io/api/rbac/v1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

var log = logf.Log.WithName("controller_workshop")

//Add - TODO add
func Add(m manager.Manager) error {
	return add(m, newReconciler(m))
}

func newReconciler(m manager.Manager) reconcile.Reconciler {
	return &ReconcileWorkshop{client: m.GetClient(), scheme: m.GetScheme()}
}

func add(m manager.Manager, r reconcile.Reconciler) error {

	//Create the workshop controller
	c, err := controller.New("workshop-controller", m, controller.Options{Reconciler: r})

	if err != nil {
		return err
	}

	//OpenShift User and Project

	if err := userv1.AddToScheme(m.GetScheme()); err != nil {
		log.Error(err, "Error adding OpenShift User API to scheme")
	}

	if err := projectv1.AddToScheme(m.GetScheme()); err != nil {
		log.Error(err, "Error adding OpenShift Project API to scheme")
	}

	if err := oauthv1.AddToScheme(m.GetScheme()); err != nil {
		log.Error(err, "Error adding OpenShift OAuth API to scheme")
	}

	// register RBAC in the scheme
	if err := rbac.AddToScheme(m.GetScheme()); err != nil {
		log.Error(err, "Failed to add RBAC to scheme")
	}
	//Watch for changes on Primary resource Workshop
	err = c.Watch(&source.Kind{Type: &workshopv1alpha1.Workshop{}}, &handler.EnqueueRequestForObject{})

	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &workshopv1alpha1.Workshop{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &workshopv1alpha1.Workshop{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.Secret{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &workshopv1alpha1.Workshop{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &rbac.Role{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &workshopv1alpha1.Workshop{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &rbac.ClusterRole{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &workshopv1alpha1.Workshop{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &rbac.RoleBinding{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &workshopv1alpha1.Workshop{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &rbac.ClusterRoleBinding{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &workshopv1alpha1.Workshop{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &oauthv1.OAuth{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &workshopv1alpha1.Workshop{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &userv1.User{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &workshopv1alpha1.Workshop{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &userv1.Group{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &workshopv1alpha1.Workshop{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &projectv1.Project{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &workshopv1alpha1.Workshop{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileWorkshop{}

//ReconcileWorkshop TODO
type ReconcileWorkshop struct {
	client client.Client
	scheme *runtime.Scheme
}

//Reconcile - the reconcile method
func (r *ReconcileWorkshop) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Workshop.")

	//Fetch Workshop Spec
	workshop := &workshopv1alpha1.Workshop{}

	err := r.client.Get(context.TODO(), request.NamespacedName, workshop)

	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			reqLogger.Info("Workshop resource not found. Ignoring since object must be deleted.")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		reqLogger.Error(err, "Failed to get Workshop.")
		return reconcile.Result{}, err
	}

	//Create workshop projects
	if workshop.Spec.Project.Create {
		reqLogger.Info("Creating Workshop Projects")
		projects := create.WorkshopProjects(workshop.Spec)
		for _, p := range projects {
			err := r.client.Create(context.TODO(), &p)
			if err != nil {
				return reconcile.Result{Requeue: true}, err
			}
			controllerutil.SetControllerReference(workshop, &p, r.scheme)
		}

	}

	//Create workshop users
	if workshop.Spec.User.Create {
		reqLogger.Info("Creating Workshop Users")
		oauth, userSecret, err := create.WorkshopUsers(workshop.Spec)
		if err != nil {
			return reconcile.Result{Requeue: true}, err
		}

		err = r.client.Get(context.TODO(),
			types.NamespacedName{Name: create.HtpassSecretNamespace, Namespace: create.HtpassSecretName}, userSecret)

		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("Creating htpass secret")
			err = r.client.Create(context.TODO(), userSecret)
		} else {
			reqLogger.Info("Updating htpass secret")
			err = r.client.Update(context.TODO(), userSecret)
		}

		if err != nil {
			return reconcile.Result{Requeue: true}, err
		}

		controllerutil.SetControllerReference(workshop, userSecret, r.scheme)

		err = r.client.Update(context.TODO(), oauth)

		if err != nil {
			return reconcile.Result{Requeue: true}, err
		}

		controllerutil.SetControllerReference(workshop, oauth, r.scheme)
	}

	//Workshop Student Group

	workshopStudentGroup := create.WorkshopStudentGroup(workshop.Spec)

	err = r.client.Get(context.TODO(),
		types.NamespacedName{Name: create.GroupName}, workshopStudentGroup)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating Workshop Students Group")
		err = r.client.Create(context.TODO(), workshopStudentGroup)
	} else {
		reqLogger.Info("Updating Workshop Students Group")
		err = r.client.Update(context.TODO(), workshopStudentGroup)
	}

	if err != nil {
		return reconcile.Result{Requeue: true}, err
	}

	controllerutil.SetControllerReference(workshop, workshopStudentGroup, r.scheme)

	//Workshop Student Role
	workshopStudentRole := create.WorkshopStudentRole()

	err = r.client.Get(context.TODO(),
		types.NamespacedName{Name: create.RoleName}, workshopStudentRole)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating Workshop Students Role")
		err = r.client.Create(context.TODO(), workshopStudentRole)
	}

	if err != nil {
		return reconcile.Result{Requeue: true}, err
	}

	controllerutil.SetControllerReference(workshop, workshopStudentRole, r.scheme)

	//Workshop Student RoleBinding
	workshopStudentRoleB := create.WorkshopStudentRoleBinding()

	err = r.client.Get(context.TODO(),
		types.NamespacedName{Name: create.RoleName}, workshopStudentRoleB)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating Workshop Students Role")
		err = r.client.Create(context.TODO(), workshopStudentRoleB)
	}

	if err != nil {
		return reconcile.Result{Requeue: true}, err
	}

	controllerutil.SetControllerReference(workshop, workshopStudentRoleB, r.scheme)

	//Create the CatalogSource Config
	stackCSCs := create.WorkshopOperatorsCatalog(workshop.Spec)

	for _, csc := range stackCSCs {
		err = r.client.Get(context.TODO(),
			types.NamespacedName{Name: csc.Name}, csc)

		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info(fmt.Sprintf("Creating CatalogSourceConfig : %s", csc.Name))
			err = r.client.Create(context.TODO(), csc)
		} else {
			reqLogger.Info(fmt.Sprintf("Creating CatalogSourceConfig : %s", csc.Name))
			err = r.client.Create(context.TODO(), csc)
		}

		if err != nil {
			return reconcile.Result{Requeue: true}, err
		}
		controllerutil.SetControllerReference(workshop, csc, r.scheme)
	}

	return reconcile.Result{}, nil
}
