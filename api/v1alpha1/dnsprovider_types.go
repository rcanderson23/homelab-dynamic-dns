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
	dnsp "github.com/rcanderson23/homelab-dynamic-dns/networking/dnsproviders"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DNSProviderSpec defines the desired state of DNSProvider
type DNSProviderSpec struct {
	Type   string      `json:"type"`
	Config dnsp.Config `json:"config"`
}

// DNSProviderStatus defines the observed state of DNSProvider
type DNSProviderStatus struct{}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=dnsproviders,scope=Cluster
// +kubebuilder:subresource:status

// DNSProvider is the Schema for the resolvers API
type DNSProvider struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DNSProviderSpec   `json:"spec,omitempty"`
	Status DNSProviderStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DNSProviderList contains a list of DNSProvider
type DNSProviderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DNSProvider `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DNSProvider{}, &DNSProviderList{})
}
