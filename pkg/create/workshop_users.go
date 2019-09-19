package create

import (
	"fmt"

	workshopv1alpha1 "github.com/kameshsampath/workshop-operator/pkg/apis/kameshs/v1alpha1"
	util "github.com/kameshsampath/workshop-operator/pkg/util"
	"golang.org/x/crypto/bcrypt"

	oauthv1 "github.com/openshift/api/config/v1"
	corev1 "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	//HtpassSecretName  the secret used for htpasswd idp
	HtpassSecretName string = "htpass-rhd-workshop"
	//HtpassSecretNamespace - the namesapce where HtpassSecretName will be created
	HtpassSecretNamespace string = "openshift-config"
	//OAuthIdentityProviderName  the OAuth Provider name
	OAuthIdentityProviderName string = "htpasswd"
	//OcpAdminRoleName is the role name that will be used for cluster-admin ocpadmin
	OcpAdminRoleName string = "ocpadmin-cluster-admin"
)

//WorkshopUsers - creates the workshop users in OpenShift
func WorkshopUsers(spec workshopv1alpha1.WorkshopSpec) ([]oauthv1.IdentityProvider, *corev1.Secret, error) {

	userHashes, err := generateUserHashes(spec.User)

	if err != nil {
		return nil, nil, err
	}

	htpassSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      HtpassSecretName,
			Namespace: HtpassSecretNamespace,
			Labels:    util.WorkshopLabels(),
		},
		Data: map[string][]byte{
			"htpasswd": []byte(userHashes),
		},
	}

	idps := []oauthv1.IdentityProvider{
		oauthv1.IdentityProvider{
			Name:          OAuthIdentityProviderName,
			MappingMethod: oauthv1.MappingMethodAdd,
			IdentityProviderConfig: oauthv1.IdentityProviderConfig{
				Type: oauthv1.IdentityProviderTypeHTPasswd,
				HTPasswd: &oauthv1.HTPasswdIdentityProvider{
					FileData: oauthv1.SecretNameReference{Name: HtpassSecretName},
				},
			},
		},
	}

	return idps, htpassSecret, nil
}

//OcpAdminRoleBinding is cluster-admin who wil be created
func OcpAdminRoleBinding() *rbac.ClusterRoleBinding {
	return &rbac.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:   OcpAdminRoleName,
			Labels: util.WorkshopLabels(),
		},
		RoleRef: rbac.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Name:     "cluster-admin",
			Kind:     "ClusterRole",
		},
		Subjects: []rbac.Subject{
			rbac.Subject{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "User",
				Name:     "ocpadmin",
			},
		},
	}
}

func generateUserHashes(u workshopv1alpha1.WorkshopUser) (string, error) {
	var userHashes string

	start := u.Start
	end := u.End
	prefix := u.Prefix
	userPassword := u.Password
	ocpAdminPassword := u.AdminPassword

	for i := start; i <= end; i++ {
		pb := []byte(userPassword)
		hb, err := bcrypt.GenerateFromPassword(pb, bcrypt.MinCost)
		if err != nil {
			return userHashes, err
		}
		userHashes = userHashes + fmt.Sprintf("%s%02d:%s\n", prefix, i, string(hb))
	}

	//add another admin user
	pb := []byte(ocpAdminPassword)
	hb, err := bcrypt.GenerateFromPassword(pb, bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	userHashes = userHashes + fmt.Sprintf("ocpadmin:%s\n", string(hb))

	return userHashes, nil
}
