package workshop

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/kameshsampath/workshop-operator/pkg/create"
	"k8s.io/apimachinery/pkg/api/errors"

	workshopv1alpha1 "github.com/kameshsampath/workshop-operator/pkg/apis/kameshs/v1alpha1"

	oauthv1 "github.com/openshift/api/config/v1"

	corev1 "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileWorkshop) createUsersRolesAndRoleBindings(reqLogger logr.Logger, workshop *workshopv1alpha1.Workshop) error {
	reqLogger.Info("Creating Workshop Users")

	isUsersUpdate := false

	idps, userSecret, _ := create.WorkshopUsers(workshop.Spec)

	us := &corev1.Secret{}

	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      create.HtpassSecretName,
		Namespace: create.HtpassSecretNamespace,
	}, us)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating Users secret")
		err = r.client.Create(context.TODO(), userSecret)
		if err != nil {
			return err
		}
		controllerutil.SetControllerReference(workshop, userSecret, r.scheme)
		isUsersUpdate = true
	} else {
		reqLogger.Info("Updating Users secret")
		us.Data["htpasswd"] = userSecret.Data["htpasswd"]
		err = r.client.Update(context.TODO(), us)
		if err != nil {
			return err
		}
		isUsersUpdate = true
	}

	if isUsersUpdate {
		oauth := &oauthv1.OAuth{}
		err = r.client.Get(context.TODO(),
			types.NamespacedName{Name: "cluster"}, oauth)

		if err != nil {
			return err
		}

		reqLogger.Info("Updating OAuth ")
		oauth.Spec.IdentityProviders = idps
		err = r.client.Update(context.TODO(), oauth)
		if err != nil {
			return err
		}
	}

	//OCP Admin
	ocpadminRoleB := create.OcpAdminRoleBinding()

	err = r.client.Get(context.TODO(),
		types.NamespacedName{Name: create.OcpAdminRoleName}, ocpadminRoleB)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating Cluster Admin Role")
		err = r.client.Create(context.TODO(), ocpadminRoleB)
	}

	if err != nil {
		return err
	}

	controllerutil.SetControllerReference(workshop, ocpadminRoleB, r.scheme)

	//Workshop Student Group

	workshopStudentGroup := create.WorkshopStudentGroup(workshop.Spec)

	err = r.client.Get(context.TODO(),
		types.NamespacedName{Name: create.GroupName}, workshopStudentGroup)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating Workshop Students Group")
		err = r.client.Create(context.TODO(), workshopStudentGroup)
		if err != nil {
			return err
		}
		controllerutil.SetControllerReference(workshop, workshopStudentGroup, r.scheme)
	} else {
		reqLogger.Info("Updating Workshop Students Group")
		err = r.client.Update(context.TODO(), workshopStudentGroup)
		if err != nil {
			return err
		}
	}

	//Workshop Student Role
	workshopStudentRole := create.WorkshopStudentRole()

	wr := &rbac.ClusterRole{}

	err = r.client.Get(context.TODO(),
		types.NamespacedName{Name: create.RoleName}, wr)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating Workshop Students Role")
		err = r.client.Create(context.TODO(), workshopStudentRole)
		if err != nil {
			return err
		}
		controllerutil.SetControllerReference(workshop, workshopStudentRole, r.scheme)
	} else {
		reqLogger.Info("Update Workshop Students Role")
		wr.Rules = workshopStudentRole.Rules
		err = r.client.Update(context.TODO(), wr)
		if err != nil {
			return err
		}
	}

	//Workshop Student RoleBinding
	workshopStudentRoleB := create.WorkshopStudentRoleBinding()

	wrb := &rbac.ClusterRoleBinding{}

	err = r.client.Get(context.TODO(),
		types.NamespacedName{Name: create.RoleName}, wrb)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating Workshop Students Role")
		err = r.client.Create(context.TODO(), workshopStudentRoleB)
	}

	if err != nil {
		return err
	}

	controllerutil.SetControllerReference(workshop, workshopStudentRoleB, r.scheme)

	return nil
}
