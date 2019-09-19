package workshop

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	workshopv1alpha1 "github.com/kameshsampath/workshop-operator/pkg/apis/kameshs/v1alpha1"
	"github.com/kameshsampath/workshop-operator/pkg/create"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	chev1 "github.com/eclipse/che-operator/pkg/apis/org/v1"
	oauthv1 "github.com/openshift/api/config/v1"
	projectv1 "github.com/openshift/api/project/v1"
	userv1 "github.com/openshift/api/user/v1"
	olmv1 "github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators"
	marketplacev2 "github.com/operator-framework/operator-marketplace/pkg/apis/operators/v2"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
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
			Stack: workshopv1alpha1.WorkshopStack{
				Install: true,
				Che: workshopv1alpha1.EclipseChe{
					CheVersion: "7.1.0",
					TLSSupport: true,
				},
				Community: []workshopv1alpha1.PackageInfo{
					workshopv1alpha1.PackageInfo{
						Name:     "knative-serving",
						Operator: "knative-serving-operator",
						Version:  "0.7.1",
					},
					workshopv1alpha1.PackageInfo{
						Name:     "pipelines",
						Operator: "openshift-pipelines-operator",
						Version:  "0.5.2",
					},
				},
				RedHat: []workshopv1alpha1.PackageInfo{
					workshopv1alpha1.PackageInfo{
						Name:     "elastic-search",
						Operator: "elasticsearch-operator",
						Version:  "4.1.15-201909041605",
					},
					workshopv1alpha1.PackageInfo{
						Name:     "jaeger",
						Operator: "jaeger-product",
						Version:  "1.13.1",
					},
					workshopv1alpha1.PackageInfo{
						Name:     "kiali",
						Operator: "kiali-ossm",
						Version:  "1.0.5",
					},
				},
			},
		},
	}

	// Objects to track in the fake client.
	objs := []runtime.Object{
		workshop,
	}

	oauth := &oauthv1.OAuth{}
	secret := &corev1.Secret{}
	project := &projectv1.Project{}
	clusterRole := &rbac.ClusterRole{}
	clusterRoleBinding := &rbac.ClusterRoleBinding{}
	userGroup := &userv1.Group{}
	csc := &marketplacev2.CatalogSourceConfig{}
	og := &olmv1.OperatorGroup{}
	sub := &olmv1.Subscription{}
	cheC := &chev1.CheCluster{}

	// Register operator types with the runtime scheme.
	s := scheme.Scheme
	s.AddKnownTypes(workshopv1alpha1.SchemeGroupVersion, workshop)
	s.AddKnownTypes(projectv1.SchemeGroupVersion, project)
	s.AddKnownTypes(oauthv1.SchemeGroupVersion, oauth)
	s.AddKnownTypes(corev1.SchemeGroupVersion, secret)
	s.AddKnownTypes(rbac.SchemeGroupVersion, clusterRole)
	s.AddKnownTypes(rbac.SchemeGroupVersion, clusterRoleBinding)
	s.AddKnownTypes(userv1.SchemeGroupVersion, userGroup)
	s.AddKnownTypes(marketplacev2.SchemeGroupVersion, csc)
	s.AddKnownTypes(olmv1.SchemeGroupVersion, og)
	s.AddKnownTypes(olmv1.SchemeGroupVersion, sub)
	s.AddKnownTypes(chev1.SchemeGroupVersion, cheC)

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
	for i := workshop.Spec.User.Start; i <= workshop.Spec.User.End; i++ {
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

	expectedIdpName := create.OAuthIdentityProviderName
	expectedIdpFileData := create.HtpassSecretName

	if actualIdpName := idp.Name; actualIdpName != expectedIdpName {
		t.Errorf("Expected IDP Name: %s but got: %s", expectedIdpName, actualIdpName)
	}
	if actualIdpFileData := idp.IdentityProviderConfig.HTPasswd.FileData.Name; actualIdpFileData != expectedIdpFileData {
		t.Errorf("Expected IDP Provider Data Name: %s but got: %s", expectedIdpName, actualIdpFileData)
	}

	ug := &userv1.Group{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: create.GroupName}, ug)

	if err != nil {
		t.Fatalf("get user group: (%v)", err)
	}
	var eug = []string{}

	for i := workshop.Spec.User.Start; i <= workshop.Spec.User.End; i++ {
		userName := fmt.Sprintf("%s%02d", workshop.Spec.User.Prefix, i)
		eug = append(eug, userName)
	}
	expectedUsersInGroup := userv1.OptionalNames(eug)
	actualUsersInGroup := ug.Users

	if !reflect.DeepEqual(expectedUsersInGroup, actualUsersInGroup) {
		t.Errorf("Expect Users in Group: (%v) but got : (%v)  ", expectedUsersInGroup, actualUsersInGroup)
	}

	cr := &rbac.ClusterRole{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: create.RoleName}, cr)

	if err != nil {
		t.Fatalf("get cluster role: (%v)", err)
	}

	crb := &rbac.ClusterRoleBinding{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: create.RoleName}, crb)

	if err != nil {
		t.Fatalf("get cluster role binding: (%v)", err)
	}

	rhcsc := &marketplacev2.CatalogSourceConfig{}
	err = cl.Get(context.TODO(), types.NamespacedName{Namespace: "openshift-marketplace", Name: create.RHPackagesCSC}, rhcsc)

	if err != nil {
		t.Fatalf("get rh catalogsource configs: (%v)", err)
	}

	expectedPackages := "elasticsearch-operator,jaeger-product,kiali-ossm"
	if pkgNames := rhcsc.Spec.Packages; expectedPackages != pkgNames {
		t.Errorf("RH Expecting packages : (%s) but got (%s) ", pkgNames, expectedPackages)
	}

	ccsc := &marketplacev2.CatalogSourceConfig{}
	err = cl.Get(context.TODO(), types.NamespacedName{Namespace: "openshift-marketplace", Name: create.CommunityPackagesCSC}, ccsc)

	if err != nil {
		t.Fatalf("get community catalogsource configs: (%v)", err)
	}

	expectedPackages = "knative-serving-operator,openshift-pipelines-operator"
	if pkgNames := ccsc.Spec.Packages; expectedPackages != pkgNames {
		t.Errorf("Community Expecting packages : (%s) but got (%s) ", pkgNames, expectedPackages)
	}

	// Nexus Checks
	nexus := &appsv1.Deployment{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: create.NexusDeploymentName, Namespace: create.RHDWorkshopNamespace}, nexus)
	if err != nil {
		t.Fatalf("get nexus deployment: (%v)", err)
	}

	if actualImage := nexus.Spec.Template.Spec.Containers[0].Image; create.NexusImageName != actualImage {
		t.Errorf("Expecting Nexus Image to be (%s) but got %s ", create.NexusImageName, actualImage)
	}

	nexusSvc := &corev1.Service{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: create.NexusDeploymentName, Namespace: create.RHDWorkshopNamespace}, nexusSvc)
	if err != nil {
		t.Fatalf("get nexus service: (%v)", err)
	}

	if actualNexusSvcPort := nexusSvc.Spec.Ports[0].Port; 8081 != actualNexusSvcPort {
		t.Errorf("Expecting Nexus Service to be 8081 but got %d ", actualNexusSvcPort)
	}

	//Che checks
	cheProject := &projectv1.Project{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: create.CheInstallNamespace}, cheProject)
	if err != nil {
		t.Fatalf("get che project: (%v)", err)
	}

	cheCSC := &marketplacev2.CatalogSourceConfig{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: create.CheCSCName, Namespace: "openshift-marketplace"}, cheCSC)
	if err != nil {
		t.Fatalf("get che catalog source config name: (%v)", cheCSC)
	}

	cheOG := &olmv1.OperatorGroup{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: create.CheOperatorGroupName, Namespace: create.CheInstallNamespace}, cheOG)
	if err != nil {
		t.Fatalf("get che operator group name: (%v)", err)
	}

	cheSub := &olmv1.Subscription{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: create.CheSubscriptionName, Namespace: create.CheInstallNamespace}, cheSub)
	if err != nil {
		t.Fatalf("get che subscription  name: (%v)", err)
	}

	cheCluster := &chev1.CheCluster{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: create.CheSubscriptionName, Namespace: create.CheInstallNamespace}, cheCluster)
	if err != nil {
		t.Fatalf("get che cluster name: (%v)", err)
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
