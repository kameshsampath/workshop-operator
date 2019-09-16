package workshop

import (
	"context"
	"fmt"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	workshopv1alpha1 "github.com/kameshsampath/workshop-operator/pkg/apis/kameshs/v1alpha1"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	oauthv1 "github.com/openshift/api/config/v1"
	projectv1 "github.com/openshift/api/project/v1"
	corev1 "k8s.io/api/core/v1"
)

func TestWorkshopController(t *testing.T) {

	logf.SetLogger(logf.ZapLogger(true))

	var (
		name      = "workshop-operator"
		namespace = "rhd-workshop-infra"
	)

	workshop := &workshopv1alpha1.Workshop{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: workshopv1alpha1.WorkshopSpec{
			Clean:              false,
			OpenshiftAPIServer: "",
			Project: workshopv1alpha1.WorkshopProject{
				Create:   true,
				Prefixes: []string{"tutorial"},
			},
			User: workshopv1alpha1.WorkshopUser{
				Create:        true,
				Start:         0,
				End:           5,
				Prefix:        "user",
				Password:      "openshift",
				AdminPassword: "openshift",
			},
			Stack: workshopv1alpha1.WorkshopStack{},
		},
	}

	// Objects to track in the fake client.
	objs := []runtime.Object{
		workshop,
	}

	oauth := &oauthv1.OAuth{}
	secret := &corev1.Secret{}
	project := &projectv1.Project{}

	// Register operator types with the runtime scheme.
	s := scheme.Scheme
	s.AddKnownTypes(workshopv1alpha1.SchemeGroupVersion, workshop)
	s.AddKnownTypes(projectv1.SchemeGroupVersion, project)
	s.AddKnownTypes(oauthv1.SchemeGroupVersion, oauth)
	s.AddKnownTypes(corev1.SchemeGroupVersion, secret)

	// Create a fake client to mock API calls.
	cl := fake.NewFakeClient(objs...)
	defaultOAuth := &oauthv1.OAuth{
		ObjectMeta: metav1.ObjectMeta{
			Name: "cluster",
		},
		Spec: oauthv1.OAuthSpec{IdentityProviders: []oauthv1.IdentityProvider{}},
	}

	err := cl.Create(context.TODO(), defaultOAuth)
	if err != nil {
		t.Fatalf("Error intializing test (%v)", err)
	}

	// Create a ReconcileWorkshop object with the scheme and fake client.
	r := &ReconcileWorkshop{client: cl, scheme: s}

	// Mock request to simulate Reconcile() being called on an event for a
	// watched resource .
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		},
	}

	refresh(t, r, req)

	//Check project has been created
	for i := 0; i <= 5; i++ {
		projectName := fmt.Sprintf("%s-%d", workshop.Spec.Project.Prefixes[0], i)

		p := &projectv1.Project{}

		err := cl.Get(context.TODO(), types.NamespacedName{Name: projectName}, p)

		if err != nil {
			t.Fatalf("get projects: (%v)", err)
		}

		annotations := p.ObjectMeta.Annotations
		expectedDesc := fmt.Sprintf("Project %s", projectName)
		expectedRequester := fmt.Sprintf("%s%02d", workshop.Spec.User.Prefix, i)

		if actualDisplayName := annotations["openshift.io/display-name"]; actualDisplayName != projectName {
			t.Errorf("Expected Project Display Name: %s but got: %s ", actualDisplayName, projectName)
		}

		if actualDesc := annotations["openshift.io/description"]; actualDesc != expectedDesc {
			t.Errorf("Expected Project Desc: %s but got: %s ", expectedDesc, actualDesc)
		}

		if actualRequester := annotations["openshift.io/requester"]; actualRequester != expectedRequester {
			t.Errorf("Expected Project Desc: %s but got: %s ", expectedRequester, actualRequester)
		}

	}

	htps := &corev1.Secret{}

	//Check htpass-bcrypt secret has been created
	err = cl.Get(context.TODO(), types.NamespacedName{Namespace: "openshift-config", Name: "htpass-bcrypt"}, htps)

	if err != nil {
		t.Fatalf("get secrets: (%v)", err)
	}

	o := &oauthv1.OAuth{}

	err = cl.Get(context.TODO(), types.NamespacedName{Name: "cluster"}, o)

	if err != nil {
		t.Fatalf("get oauth: (%v)", err)
	}

	if len(o.Spec.IdentityProviders) != 1 {
		t.Errorf("Expect 1 oauth identity providers but got zero (%v) ", err)
	}

	idp := o.Spec.IdentityProviders[0]

	expectedIdpName := "htpasswd"
	expectedIdpFileData := "htpass-bcrypt"

	if actualIdpName := idp.Name; actualIdpName != expectedIdpName {
		t.Errorf("Expected IDP Name: %s but got: %s", expectedIdpName, actualIdpName)
	}
	if actualIdpFileData := idp.IdentityProviderConfig.HTPasswd.FileData.Name; actualIdpFileData != expectedIdpFileData {
		t.Errorf("Expected IDP Provider Data Name: %s but got: %s", expectedIdpName, actualIdpFileData)
	}
}

func refresh(t *testing.T, r *ReconcileWorkshop, req reconcile.Request) {
	res, err := r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}
	// Check the result of reconciliation to make sure it has the desired state.
	if res.Requeue {
		t.Error("reconcile requeue which is not expected")
	}
}
