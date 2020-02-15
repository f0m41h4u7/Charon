package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CharonSpec struct {
	Registry string `json:"registry"`
        Image string `json:"image"`
        Version string `json:"version"`
}

type CharonStatus struct {
        Registry string `json:"registry"`
        Image string `json:"image"`
        Version string `json:"version"`
        VersionChanged bool `json:"version_changed"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:subresource:status
// +kubebuilder:resource:path=charons,scope=Namespaced
type Charon struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CharonSpec   `json:"spec,omitempty"`
	Status CharonStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type CharonList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Charon `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Charon{}, &CharonList{})
}
