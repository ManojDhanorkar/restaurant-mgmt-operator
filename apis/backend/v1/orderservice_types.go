/*
Copyright 2023.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// OrderServiceSpec defines the desired state of OrderService
type OrderServiceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// This spec allows you to supply image name
	Image string `json:"image"`

	// This spec allows you to supply size of deployment
	Size int32 `json:"size"`
}

// OrderServiceStatus defines the observed state of OrderService
type OrderServiceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// Nodes are the names of the order service pods
	Nodes []string `json:"nodes"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// OrderService is the Schema for the orderservices API
type OrderService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OrderServiceSpec   `json:"spec,omitempty"`
	Status OrderServiceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// OrderServiceList contains a list of OrderService
type OrderServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OrderService `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OrderService{}, &OrderServiceList{})
}
