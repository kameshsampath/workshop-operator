package create

import (
	"fmt"

	workshopv1alpha1 "github.com/kameshsampath/workshop-operator/pkg/apis/kameshs/v1alpha1"
	projectv1 "github.com/openshift/api/project/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//WorkshopProjects - creates the workshop projects in OpenShift
func WorkshopProjects(spec workshopv1alpha1.WorkshopSpec) []projectv1.Project {

	var projects []projectv1.Project
	start := spec.User.Start
	end := spec.User.End
	prefixes := spec.Project.Prefixes

	for i := start; i <= end; i++ {
		for _, name := range prefixes {
			pName := fmt.Sprintf("%s-%d", name, i)
			pDesc := fmt.Sprintf("Project %s", pName)
			pRequester := fmt.Sprintf("%s%02d", spec.User.Prefix, i)
			p := &projectv1.Project{
				ObjectMeta: metav1.ObjectMeta{
					Name: pName,
					Annotations: map[string]string{
						"openshift.io/display-name": pName,
						"openshift.io/description":  pDesc,
						"openshift.io/requester":    pRequester,
					},
				},
			}

			projects = append(projects, *p)
		}
	}

	return projects
}
