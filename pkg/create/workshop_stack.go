package create

import (
	"strings"

	workshopv1alpha1 "github.com/kameshsampath/workshop-operator/pkg/apis/kameshs/v1alpha1"
	marketplacev2 "github.com/operator-framework/operator-marketplace/pkg/apis/operators/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	//CSCNS the namespace where Catalog Source Config will be created
	CSCNS string = "openshift-marketplace"
	//RHPackagesCSC the Red Hat Operators CSC
	RHPackagesCSC string = "rhd-workshop-packages"
	//CommunityPackagesCSC the Community Operators CSC
	CommunityPackagesCSC string = "rhd-workshop-packages-community"
)

//WorkshopOperatorsCatalog creates the OperatorsCatalog for community-operators in workshop stack
func WorkshopOperatorsCatalog(spec workshopv1alpha1.WorkshopSpec) []*marketplacev2.CatalogSourceConfig {
	var cscs = []*marketplacev2.CatalogSourceConfig{}

	//RedHat Operators
	var redhatOperatorPkgs []string

	for _, p := range spec.Stack.RedHat {
		redhatOperatorPkgs = append(redhatOperatorPkgs, p.Operator)
	}
	cscs = append(cscs, &marketplacev2.CatalogSourceConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      RHPackagesCSC,
			Namespace: CSCNS,
		},
		Spec: marketplacev2.CatalogSourceConfigSpec{
			TargetNamespace: "openshift-operators",
			Source:          "redhat-operators",
			Packages:        strings.Join(redhatOperatorPkgs, ","),
		},
	})

	var communityOperatorPkgs []string
	for _, p := range spec.Stack.Community {
		communityOperatorPkgs = append(communityOperatorPkgs, p.Operator)
	}

	//Community Operators
	cscs = append(cscs, &marketplacev2.CatalogSourceConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      CommunityPackagesCSC,
			Namespace: CSCNS,
		},
		Spec: marketplacev2.CatalogSourceConfigSpec{
			TargetNamespace: "openshift-operators",
			Source:          "community-operators",
			Packages:        strings.Join(communityOperatorPkgs, ","),
		},
	})

	return cscs
}
