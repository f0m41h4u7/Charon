package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CharonSpec struct {
}

type CharonStatus struct {
}

type Charon struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CharonSpec   `json:"spec,omitempty"`
	Status CharonStatus `json:"status,omitempty"`
}

type CharonList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Charon `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Charon{}, &CharonList{})
}
