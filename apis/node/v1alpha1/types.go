package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type AuthorizedKey struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metdata,omitempty"`

	Data AuthorizedKeyData `json:"data"`
}

type AuthorizedKeyData struct {
	Key string `json:"key"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type AuthorizedKeyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metdata,omitempty"`
	Items           []AuthorizedKey `json:"items"`
}
