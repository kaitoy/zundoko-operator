/*
Copyright 2019 kaitoy.

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

package v1beta1

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// HikawaSpec defines the desired state of Hikawa
type HikawaSpec struct {
	IntervalMillis time.Duration `json:"intervalMillis"`
	NumZundokos    int           `json:"numZundokos,omitempty"`
	SayKiyoshi     bool          `json:"sayKiyoshi,omitempty"`
}

// HikawaStatus defines the observed state of Hikawa
type HikawaStatus struct {
	NumZundokosSaid int  `json:"numZundokosSaid"`
	Kiyoshied       bool `json:"kiyoshied"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Hikawa is the Schema for the hikawas API
// +k8s:openapi-gen=true
type Hikawa struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HikawaSpec   `json:"spec,omitempty"`
	Status HikawaStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HikawaList contains a list of Hikawa
type HikawaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Hikawa `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Hikawa{}, &HikawaList{})
}
