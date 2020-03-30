package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeployerSpec struct {
	DeployerImage string `json:"deployerImage"`
	AnalyzerImage string `json:"analyzerImage"`
}

type DeployerStatus struct {
	Image string `json:"image"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:subresource:status
// +kubebuilder:resource:path=deployers,scope=Namespaced
type Deployer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DeployerSpec   `json:"spec,omitempty"`
	Status DeployerStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type DeployerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Deployer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Deployer{}, &DeployerList{})
}
