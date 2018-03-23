package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AuthorizedKey defines the base parsing
type AuthorizedKey struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Data              AuthorizedKeyData `json:"data"`
}

// AuthorizedKeyData defines the spec tag parsing
type AuthorizedKeyData struct {
	Key            string `json:"key"`
	GithubUsername string `json:"githubUsername"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AuthorizedKeyList defines the collection AuthorizedKeys
type AuthorizedKeyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []AuthorizedKey `json:"items"`
}
