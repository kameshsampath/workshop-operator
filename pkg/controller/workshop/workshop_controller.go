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

	// err = c.Watch(&source.Kind{Type: &marketplacev2.CatalogSourceConfig{}}, &handler.EnqueueRequestForOwner{
	// 	IsController: true,
	// 	OwnerType:    &workshopv1alpha1.Workshop{},
	// })
	// if err != nil {
	// 	return err
	// }

	// err = c.Watch(&source.Kind{Type: &olmv1.Subscription{}}, &handler.EnqueueRequestForOwner{
	// 	IsController: true,
	// 	OwnerType:    &workshopv1alpha1.Workshop{},
	// })
	// if err != nil {
	// 	return err
	// }

	// err = c.Watch(&source.Kind{Type: &olmv1.OperatorGroup{}}, &handler.EnqueueRequestForOwner{
	// 	IsController: true,
	// 	OwnerType:    &workshopv1alpha1.Workshop{},
	// })
	// if err != nil {
	// 	return err
	// }

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
		projects := create.WorkshopProjects(workshop.Spec)
		for _, p := range projects {
			err = r.client.Get(context.TODO(),
				types.NamespacedName{Name: p.Name}, p)
			if err != nil {
				if errors.IsNotFound(err) {
					reqLogger.Info("Creating Workshop Project", "Project Name", p.Name)
					err := r.client.Create(context.TODO(), p)
					if err != nil {
						return reconcile.Result{Requeue: true}, err
					}
					controllerutil.SetControllerReference(workshop, p, r.scheme)
				} else {
					reqLogger.Info("Creating Workshop Project", "Project Name", p.Name, "Already Exists")
				}
			}
		}

	}

	//Create workshop users
	if workshop.Spec.User.Create {
		log.Info("Creating Workshop Users")

		isUsersUpdate := false

		idps, userSecret, _ := create.WorkshopUsers(workshop.Spec)

		us := &corev1.Secret{}

		err = r.client.Get(context.TODO(), types.NamespacedName{
			Name:      create.HtpassSecretName,
			Namespace: create.HtpassSecretNamespace,
		}, us)

		log.Info("Secret:::", "", us, "Err", err)

		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("Creating Users secret")
			err = r.client.Create(context.TODO(), userSecret)
			if err != nil {
				return reconcile.Result{Requeue: true}, err
			}
			controllerutil.SetControllerReference(workshop, userSecret, r.scheme)
			isUsersUpdate = true
		} else {
			reqLogger.Info("Updating Users secret")
			us.Data["htpasswd"] = userSecret.Data["htpasswd"]
			err = r.client.Update(context.TODO(), us)
			if err != nil {
				return reconcile.Result{Requeue: true}, err
			}
			isUsersUpdate = true
		}

		if isUsersUpdate {
			oauth := &oauthv1.OAuth{}
			err = r.client.Get(context.TODO(),
				types.NamespacedName{Name: "cluster"}, oauth)

			if err != nil {
				return reconcile.Result{Requeue: true}, err
			}

			reqLogger.Info("Updating OAuth ")
			oauth.Spec.IdentityProviders = idps
			err = r.client.Update(context.TODO(), oauth)
			if err != nil {
				return reconcile.Result{Requeue: true}, err
			}
		}
		//Workshop Student Group

		workshopStudentGroup := create.WorkshopStudentGroup(workshop.Spec)

		err = r.client.Get(context.TODO(),
			types.NamespacedName{Name: create.GroupName}, workshopStudentGroup)

		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("Creating Workshop Students Group")
			err = r.client.Create(context.TODO(), workshopStudentGroup)
			if err != nil {
				return reconcile.Result{Requeue: true}, err
			}
			controllerutil.SetControllerReference(workshop, workshopStudentGroup, r.scheme)
		} else {
			reqLogger.Info("Updating Workshop Students Group")
			err = r.client.Update(context.TODO(), workshopStudentGroup)
			if err != nil {
				return reconcile.Result{Requeue: true}, err
			}
		}

		//Workshop Student Role
		workshopStudentRole := create.WorkshopStudentRole()

		err = r.client.Get(context.TODO(),
			types.NamespacedName{Name: create.RoleName}, workshopStudentRole)

		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("Creating Workshop Students Role")
			err = r.client.Create(context.TODO(), workshopStudentRole)
			if err != nil {
				return reconcile.Result{Requeue: true}, err
			}
			controllerutil.SetControllerReference(workshop, workshopStudentRole, r.scheme)
		} else {
			reqLogger.Info("Update Workshop Students Role")
			err = r.client.Update(context.TODO(), workshopStudentRole)
			if err != nil {
				return reconcile.Result{Requeue: true}, err
			}
		}

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

	}
	//Create the CatalogSource Config
	stackCSCs := create.WorkshopOperatorsCatalog(workshop.Spec)

	for _, csc := range stackCSCs {
		err = r.client.Get(context.TODO(),
			types.NamespacedName{Name: csc.Name}, csc)

		if err != nil {
			reqLogger.Info(fmt.Sprintf("Creating CatalogSourceConfig : %s", csc.Name))
			err = r.client.Create(context.TODO(), csc)
			if err != nil {
				return reconcile.Result{Requeue: true}, err
			}
			controllerutil.SetControllerReference(workshop, csc, r.scheme)
		} else {
			reqLogger.Info(fmt.Sprintf("Updating CatalogSourceConfig : %s", csc.Name))
			err = r.client.Update(context.TODO(), csc)
			if err != nil {
				return reconcile.Result{Requeue: true}, err
			}
		}

	}

	if workshop.Spec.Stack.Install {
		//Install Nexus
		nexus := create.Nexus()
		err = r.client.Get(context.TODO(),
			types.NamespacedName{Name: create.NexusDeploymentName, Namespace: create.RHDWorkshopNamespace}, nexus)

		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("Deploying Nexus")
			err = r.client.Create(context.TODO(), nexus)
			if err != nil {
				return reconcile.Result{Requeue: true}, err
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
				return reconcile.Result{Requeue: true}, err
			}
			controllerutil.SetControllerReference(workshop, nexusSvc, r.scheme)
		}

		//Install Eclipse Che
		cheProject := create.CheProject()
		err = r.client.Get(context.TODO(),
			types.NamespacedName{Name: create.CheInstallNamespace}, cheProject)

		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("Creating Che project")
			err = r.client.Create(context.TODO(), cheProject)
			if err != nil {
				return reconcile.Result{Requeue: true}, err
			}
			controllerutil.SetControllerReference(workshop, cheProject, r.scheme)
		}

		cheCSC := create.CheCatalogSourceConfig()
		err = r.client.Get(context.TODO(),
			types.NamespacedName{Name: create.CheCSCName, Namespace: "openshift-marketplace"}, cheCSC)

		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("Creating Che CatalogSource Config")
			err = r.client.Create(context.TODO(), cheCSC)
			if err != nil {
				return reconcile.Result{Requeue: true}, err
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
				return reconcile.Result{Requeue: true}, err
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
				return reconcile.Result{Requeue: true}, err
			}
			controllerutil.SetControllerReference(workshop, cheSub, r.scheme)
		}

		che := create.CheCluster(workshop.Spec)
		err = r.client.Get(context.TODO(),
			types.NamespacedName{Name: create.CheSubscriptionName, Namespace: create.CheInstallNamespace}, che)

		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("Creating Che Che Cluster")
			err = r.client.Create(context.TODO(), che)
			if err != nil {
				return reconcile.Result{Requeue: true}, err
			}
			controllerutil.SetControllerReference(workshop, che, r.scheme)
		}
	}

	return reconcile.Result{}, nil
}
