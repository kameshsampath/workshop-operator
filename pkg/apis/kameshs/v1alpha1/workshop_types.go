package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//WorkshopProject defines the student projects of Workshop
type WorkshopProject struct {
	Create   bool     `json:"create"`
	Prefixes []string `json:"prefixes"`
}

//WorkshopUser defines the students for Workshop
type WorkshopUser struct {
	Create        bool   `json:"create"`
	Prefix        string `json:"prefix"`
	Password      string `json:"password"`
	AdminPassword string `json:"adminPassword"`
	//+kubebuilder:validation:Minimum=0
	//+kubebuilder:validation:Maximum=1
	Start int32 `json:"start"`
	//+kubebuilder:validation:Minimum=1
	//+kubebuilder:validation:Maximum=100
	End int32 `json:"end"`
}

//WorkshopStack defines the software stack for the Workshop
type WorkshopStack struct {
	Install   bool          `json:"install"`
	Che       EclipseChe    `json:"che"`
	Community []PackageInfo `json:"community"`
	RedHat    []PackageInfo `json:"redhat"`
}

//PackageInfo defines the software package that will be installed using the operators
type PackageInfo struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Operator string `json:"operator"`
}

//EclipseChe configures few parameters for CR Eclipse Che Cluster
type EclipseChe struct {
	EnableTokenExchange bool   `json:"enableTokenExchange"`
	OperatorImage       string `json:"operatorImage"`
	CheVersion          string `json:"version"`
	TLSSupport          bool   `json:"tlsSupport,omitempty"`
	SelfSignedCert      bool   `json:"selfSignedCert,omitempty"`
}

// WorkshopSpec defines the desired state of Workshop
// +k8s:openapi-gen=true
type WorkshopSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Clean              bool            `json:"clean"`
	OpenshiftAPIServer string          `json:"openshiftAPIServer,omitempty"`
	Project            WorkshopProject `json:"project"`
	User               WorkshopUser    `json:"user"`
	Stack              WorkshopStack   `json:"stack"`
}

// WorkshopStatus defines the observed state of Workshop
// +k8s:openapi-gen=true
type WorkshopStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Workshop is the Schema for the workshops API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Workshop struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WorkshopSpec   `json:"spec,omitempty"`
	Status WorkshopStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// WorkshopList contains a list of Workshop
type WorkshopList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Workshop `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Workshop{}, &WorkshopList{})
}
