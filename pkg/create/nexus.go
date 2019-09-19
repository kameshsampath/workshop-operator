package create

import (
	util "github.com/kameshsampath/workshop-operator/pkg/util"
	appsv1 "github.com/openshift/api/apps/v1"
	isv1 "github.com/openshift/api/image/v1"
	corev1 "k8s.io/api/core/v1"

	//"k8s.io/apimachinery/pkg/api/resource"
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

//NexusImageStream creates a nexus image stream to use in nexus app deployment
func NexusImageStream() *isv1.ImageStream {
	return &isv1.ImageStream{
		ObjectMeta: metav1.ObjectMeta{
			Name:      NexusDeploymentName,
			Namespace: RHDWorkshopNamespace,
			Labels:    util.WorkshopLabels(),
		},
		Spec: isv1.ImageStreamSpec{
			LookupPolicy: isv1.ImageLookupPolicy{
				Local: true,
			},
			Tags: []isv1.TagReference{
				isv1.TagReference{
					Name: "latest",
					Annotations: map[string]string{
						"openshift.io/imported-from": NexusImageName,
					},
					From: &corev1.ObjectReference{
						Kind: "DockerImage",
						Name: NexusImageName,
					},
					ReferencePolicy: isv1.TagReferencePolicy{
						Type: isv1.SourceTagReferencePolicy,
					},
				},
			},
		},
	}
}

//Nexus creates Nexus Deployment
func Nexus() *appsv1.DeploymentConfig {
	return &appsv1.DeploymentConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      NexusDeploymentName,
			Namespace: RHDWorkshopNamespace,
			Labels:    util.WorkshopLabels(),
		},
		Spec: appsv1.DeploymentConfigSpec{
			Replicas: 1,
			Selector: map[string]string{
				"app": NexusDeploymentName,
			},
			Triggers: appsv1.DeploymentTriggerPolicies{
				appsv1.DeploymentTriggerPolicy{
					Type: appsv1.DeploymentTriggerOnConfigChange,
				},
				appsv1.DeploymentTriggerPolicy{
					Type: appsv1.DeploymentTriggerOnImageChange,
					ImageChangeParams: &appsv1.DeploymentTriggerImageChangeParams{
						Automatic:      true,
						ContainerNames: []string{"nexus"},
						From: corev1.ObjectReference{
							Kind:      "ImageStreamTag",
							Namespace: RHDWorkshopNamespace,
							Name:      "nexus:latest",
						},
					},
				},
			},
			Template: &corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": NexusDeploymentName,
					},
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						corev1.Volume{
							Name: "nexus-volume-1",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
					Containers: []corev1.Container{
						corev1.Container{
							Name:  NexusDeploymentName,
							Image: NexusImageName,
							VolumeMounts: []corev1.VolumeMount{
								corev1.VolumeMount{
									Name:      "nexus-volume-1",
									MountPath: "/sonatype-work",
								},
							},
							//TODO fix this
							// Resources: corev1.ResourceRequirements{
							// 	Requests: corev1.ResourceList{
							// 		corev1.ResourceMemory: resource.MustParse("512m"),
							// 		corev1.ResourceCPU:    resource.MustParse("500m"),
							// 	},
							// 	Limits: corev1.ResourceList{
							// 		corev1.ResourceMemory: resource.MustParse("1024m"),
							// 		corev1.ResourceCPU:    resource.MustParse("1Gi"),
							// 	},
							// },

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
			Labels:    util.WorkshopLabels(),
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
