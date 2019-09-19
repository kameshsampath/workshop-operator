package create

import (
	projectv1 "github.com/openshift/api/project/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	chev1 "github.com/eclipse/che-operator/pkg/apis/org/v1"
	workshopv1alpha1 "github.com/kameshsampath/workshop-operator/pkg/apis/kameshs/v1alpha1"
	olmv1 "github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators"
	marketplacev2 "github.com/operator-framework/operator-marketplace/pkg/apis/operators/v2"
)

const (
	//CheInstallNamespace is the namespace where Eclipse Che will be installed
	CheInstallNamespace string = "che"
	//CheOperatorGroupName the operator group for Eclipse Che Operator
	CheOperatorGroupName string = "che-workspaces"
	//CheCSCName is the CatalogSourceConfig  used in creating Eclipse Che Subscription
	CheCSCName string = "rhd-workshop-che"
	//CheSubscriptionName is the Eclipse Che Operator Subscription Name
	CheSubscriptionName string = "eclipse-che"
)

//CheProject Creates the Eclipse Che Project
func CheProject() *projectv1.Project {
	return &projectv1.Project{
		ObjectMeta: metav1.ObjectMeta{
			Name: CheInstallNamespace,
			Annotations: map[string]string{
				"openshift.io/display-name": "Eclipse Che",
				"openshift.io/description":  "Project where Eclipse Che components are installed",
			},
		},
	}
}

//CheCatalogSourceConfig Creates the Eclipse Che Catalog Source
func CheCatalogSourceConfig() *marketplacev2.CatalogSourceConfig {
	return &marketplacev2.CatalogSourceConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      CheCSCName,
			Namespace: "openshift-marketplace",
		},
		Spec: marketplacev2.CatalogSourceConfigSpec{
			TargetNamespace: CheInstallNamespace,
			Source:          "community-operators",
			Packages:        "eclipse-che",
		},
	}
}

//CheOperatorGroup Creates the Eclipse Che Operator Group
func CheOperatorGroup() *olmv1.OperatorGroup {
	return &olmv1.OperatorGroup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      CheOperatorGroupName,
			Namespace: CheInstallNamespace,
		},
		Spec: olmv1.OperatorGroupSpec{
			TargetNamespaces: []string{CheInstallNamespace},
		},
	}
}

//CheSubscription Creates the Eclipse Che Subscription
func CheSubscription(spec workshopv1alpha1.WorkshopSpec) *olmv1.Subscription {
	return &olmv1.Subscription{
		ObjectMeta: metav1.ObjectMeta{
			Name:      CheSubscriptionName,
			Namespace: CheInstallNamespace,
		},
		Spec: &olmv1.SubscriptionSpec{
			Channel:                "stable",
			Package:                "eclipse-che",
			StartingCSV:            "eclipse-che.v" + spec.Stack.Che.CheVersion,
			CatalogSourceNamespace: CheCSCName,
		},
	}
}

//CheCluster Creates the Eclipse Che Clusters
func CheCluster(spec workshopv1alpha1.WorkshopSpec) *chev1.CheCluster {

	return &chev1.CheCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      CheSubscriptionName,
			Namespace: CheInstallNamespace,
		},
		Spec: chev1.CheClusterSpec{
			Server: chev1.CheClusterSpecServer{
				TlsSupport:     spec.Stack.Che.TLSSupport,
				SelfSignedCert: spec.Stack.Che.SelfSignedCert,
			},
			Database: chev1.CheClusterSpecDB{
				ExternalDB: false,
			},
			Auth: chev1.CheClusterSpecAuth{
				OpenShiftOauth:   true,
				ExternalKeycloak: false,
			},
			Storage: chev1.CheClusterSpecStorage{
				PvcStrategy:       "per-workspace",
				PvcClaimSize:      "1Gi",
				PreCreateSubPaths: true,
			},
		},
	}
}
