package workshop

import (
	"context"
	"strings"
	"time"

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

	osappsv1 "github.com/openshift/api/apps/v1"
	oauthv1 "github.com/openshift/api/config/v1"
	isv1 "github.com/openshift/api/image/v1"
	projectv1 "github.com/openshift/api/project/v1"
	userv1 "github.com/openshift/api/user/v1"
	rbac "k8s.io/api/rbac/v1"

	chev1 "github.com/eclipse/che-operator/pkg/apis"
	cheorgv1 "github.com/eclipse/che-operator/pkg/apis/org/v1"
	olm "github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators"
	olmv1 "github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators/v1"
	olmv1alpha1 "github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators/v1alpha1"
	marketplace "github.com/operator-framework/operator-marketplace/pkg/apis"
	marketplacev2 "github.com/operator-framework/operator-marketplace/pkg/apis/operators/v2"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
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

	//OpenShift
	if err := userv1.AddToScheme(m.GetScheme()); err != nil {
		log.Error(err, "Error adding OpenShift User API to scheme")
	}

	if err := projectv1.AddToScheme(m.GetScheme()); err != nil {
		log.Error(err, "Error adding OpenShift Project API to scheme")
	}

	if err := oauthv1.AddToScheme(m.GetScheme()); err != nil {
		log.Error(err, "Error adding OpenShift OAuth API to scheme")
	}

	if err := osappsv1.AddToScheme(m.GetScheme()); err != nil {
		log.Error(err, "Error adding OpenShift Apps to scheme")
	}

	if err := isv1.AddToScheme(m.GetScheme()); err != nil {
		log.Error(err, "Error adding OpenShift Apps to scheme")
	}

	// register RBAC in the scheme

	if err := rbac.AddToScheme(m.GetScheme()); err != nil {
		log.Error(err, "Failed to add RBAC to scheme")
	}

	if err := marketplace.AddToScheme(m.GetScheme()); err != nil {
		log.Error(err, "Failed to add marketplace to scheme")

	}

	//Operators and OLM
	if err := olm.AddToScheme(m.GetScheme()); err != nil {
		log.Error(err, "Failed to add olm to scheme")
	}

	if err := olmv1alpha1.AddToScheme(m.GetScheme()); err != nil {
		log.Error(err, "Failed to add olmv1 to scheme")
	}

	if err := olmv1.AddToScheme(m.GetScheme()); err != nil {
		log.Error(err, "Failed to add olmv1 to scheme")
	}

	if err := chev1.AddToScheme(m.GetScheme()); err != nil {
		log.Error(err, "Failed to add Eclipse Che Cluster to scheme")
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

	err = c.Watch(&source.Kind{Type: &osappsv1.DeploymentConfig{}}, &handler.EnqueueRequestForOwner{
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

	err = c.Watch(&source.Kind{Type: &isv1.ImageStream{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &workshopv1alpha1.Workshop{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &marketplacev2.CatalogSourceConfig{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &workshopv1alpha1.Workshop{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &olmv1alpha1.Subscription{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &workshopv1alpha1.Workshop{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &olmv1.OperatorGroup{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &workshopv1alpha1.Workshop{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &cheorgv1.CheCluster{}}, &handler.EnqueueRequestForOwner{
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
		err = r.createProjects(log, workshop)
		if err != nil {
			return reconcile.Result{Requeue: true}, err
		}
	}

	//Create workshop users
	if workshop.Spec.User.Create {
		err = r.createUsersRolesAndRoleBindings(log, workshop)
		if err != nil {
			return reconcile.Result{Requeue: true}, err
		}
	}

	//Create workshop catalogs
	err = r.createCatalogSourceConfigs(reqLogger, workshop)
	if err != nil {
		return reconcile.Result{Requeue: true}, err
	}

	//Create workshop stacks
	if workshop.Spec.Stack.Install {
		err = r.installNexus(reqLogger, workshop)
		if err != nil {
			return reconcile.Result{Requeue: true}, err
		}
		err = r.installAndConfigureEclipseChe(reqLogger, workshop)
		if err != nil {
			if errors.IsNotFound(err) {
				if strings.Contains(err.Error(), "\"eclipse-che.v7.1.0\" not found") {
					reqLogger.Info("Che CSV is not yet available, will trying again")
					return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 10}, err
				}
			}
			return reconcile.Result{Requeue: true}, err

		}
	}

	return reconcile.Result{}, nil
}
