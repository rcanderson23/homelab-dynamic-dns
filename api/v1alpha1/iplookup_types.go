/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"github.com/rcanderson23/homelab-dynamic-dns/networking/ip"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IPLookupSpec defines the desired state of IPLookup
type IPLookupSpec struct {
	Type   string   `json:"type"`
	Config IPConfig `json:"config"`
}

// IPLookupStatus defines the observed state of IPLookup
type IPLookupStatus struct {
	Address string `json:"address,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=iplookups,scope=Cluster
// +kubebuilder:subresource:status

// IPLookup is the Schema for the iplookups API
type IPLookup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IPLookupSpec   `json:"spec,omitempty"`
	Status IPLookupStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// IPLookupList contains a list of IPLookup
type IPLookupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IPLookup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IPLookup{}, &IPLookupList{})
}

type IPConfig struct {
	Http       ip.HttpLookup `json:"http,omitempty"`
	EdgeRouter ip.EdgeRouter `json:"edgeRouter,omitempty"`
}
