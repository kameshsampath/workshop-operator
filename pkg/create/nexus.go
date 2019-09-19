package create

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	//NexusDeploymentName is the name of Nexus Deployment
	NexusDeploymentName string = "nexus"
	//RHDWorkshopNamespace is the namespace where nexus will be deployed
	RHDWorkshopNamespace string = "rhd-workshop-infra"
	//NexusImageName is the name of the container image that will be used for nexus container
	NexusImageName string = "sonatype/nexus"
)

//Nexus creates Nexus Deployment
func Nexus() *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      NexusDeploymentName,
			Namespace: RHDWorkshopNamespace,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": NexusDeploymentName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": NexusDeploymentName,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						corev1.Container{
							Name:  NexusDeploymentName,
							Image: NexusImageName,
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceMemory: resource.MustParse("512m"),
									corev1.ResourceCPU:    resource.MustParse("500m"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceMemory: resource.MustParse("1024m"),
									corev1.ResourceCPU:    resource.MustParse("1Gi"),
								},
							},
						},
					},
				},
			},
		},
	}
}

//NexusService creates Nexus service
func NexusService() *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      NexusDeploymentName,
			Namespace: RHDWorkshopNamespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": NexusDeploymentName,
			},
			Ports: []corev1.ServicePort{
				corev1.ServicePort{
					Name:     "8081-tcp",
					Port:     8081,
					Protocol: "TCP",
				},
			},
		},
	}
}
