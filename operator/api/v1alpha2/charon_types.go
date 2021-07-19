package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CharonSpec defines the desired state of Charon
type CharonSpec struct {
	Analyzer      string `json:"analyzer"`
	DeployerImage string `json:"deployerImage"`
	AnalyzerImage string `json:"analyzerImage"`
}

// CharonStatus defines the observed state of Charon
type CharonStatus struct {
	Image string `json:"image"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Charon is the Schema for the charons API
type Charon struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CharonSpec   `json:"spec,omitempty"`
	Status CharonStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CharonList contains a list of Charon
type CharonList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Charon `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Charon{}, &CharonList{})
}
