package create

import (
	"fmt"

	workshopv1alpha1 "github.com/kameshsampath/workshop-operator/pkg/apis/kameshs/v1alpha1"
	util "github.com/kameshsampath/workshop-operator/pkg/util"

	rbac "k8s.io/api/rbac/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	userv1 "github.com/openshift/api/user/v1"
)

const (
	//RoleName is the name of the role
	RoleName string = "workshop-student"
	//GroupName is the name of the groups
	GroupName string = "workshop-students"
)

//WorkshopStudentGroup - the group that represents the Workshop students
func WorkshopStudentGroup(spec workshopv1alpha1.WorkshopSpec) *userv1.Group {
	start := spec.User.Start
	end := spec.User.End
	prefix := spec.User.Prefix

	var groupUsers = []string{}

	for i := start; i <= end; i++ {
		userName := fmt.Sprintf("%s%02d", prefix, i)
		groupUsers = append(groupUsers, userName)
	}

	var group = &userv1.Group{
		ObjectMeta: metav1.ObjectMeta{
			Name:   GroupName,
			Labels: util.WorkshopLabels(),
		},
		Users: userv1.OptionalNames(groupUsers),
	}

	return group
}

//WorkshopStudentRole - the group that represents the Workshop Student Role
func WorkshopStudentRole() *rbac.ClusterRole {
	return &rbac.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:   RoleName,
			Labels: util.WorkshopLabels(),
		},
		Rules: []rbac.PolicyRule{
			rbac.PolicyRule{
				APIGroups: []string{rbac.APIGroupAll},
				Resources: []string{"configmaps", "services"},
				Verbs:     []string{"get", "view"},
			},
		},
	}
}

//WorkshopStudentRoleBinding binds wht WorkshopStudentRole with WorkshopStudent group
func WorkshopStudentRoleBinding() *rbac.ClusterRoleBinding {
	return &rbac.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:   RoleName,
			Labels: util.WorkshopLabels(),
		},
		RoleRef: rbac.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Name:     RoleName,
			Kind:     "ClusterRole",
		},
		Subjects: []rbac.Subject{
			rbac.Subject{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "Group",
				Name:     GroupName,
			},
		},
	}
}
