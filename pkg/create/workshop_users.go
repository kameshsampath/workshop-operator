package create

import (
	"fmt"

	workshopv1alpha1 "github.com/kameshsampath/workshop-operator/pkg/apis/kameshs/v1alpha1"
	"golang.org/x/crypto/bcrypt"

	oauthv1 "github.com/openshift/api/config/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	HtpassSecretName          = "htpass-bcrypt"
	HtpassSecretNamespace     = "openshift-config"
	oAuthIdentityProviderName = "htpasswd"
)

//WorkshopUsers - creates the workshop users in OpenShift
func WorkshopUsers(spec workshopv1alpha1.WorkshopSpec) (*oauthv1.OAuth, *corev1.Secret, error) {

	userHashes, err := generateUserHashes(spec.User)

	if err != nil {
		return nil, nil, err
	}

	htpassSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      HtpassSecretName,
			Namespace: HtpassSecretNamespace,
		},
		Data: map[string][]byte{
			"htpasswd": []byte(userHashes),
		},
	}

	idps := []oauthv1.IdentityProvider{
		oauthv1.IdentityProvider{
			Name:          oAuthIdentityProviderName,
			MappingMethod: oauthv1.MappingMethodAdd,
			IdentityProviderConfig: oauthv1.IdentityProviderConfig{
				Type: oauthv1.IdentityProviderTypeHTPasswd,
				HTPasswd: &oauthv1.HTPasswdIdentityProvider{
					FileData: oauthv1.SecretNameReference{Name: htpassSecret.ObjectMeta.Name},
				},
			},
		},
	}

	htpasswdIdp := &oauthv1.OAuth{
		ObjectMeta: metav1.ObjectMeta{Name: "cluster"},
		Spec: oauthv1.OAuthSpec{
			IdentityProviders: idps,
		},
	}

	return htpasswdIdp, htpassSecret, nil
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
