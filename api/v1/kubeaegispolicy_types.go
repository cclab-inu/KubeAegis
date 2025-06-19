/*
Copyright 2025.

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

// KubeAegisPolicySpec defines the desired state of KubeAegisPolicy.
type KubeAegisPolicySpec struct {
	EnableReporting bool            `json:"enableReport,omitempty"`
	IntentRequest   []IntentRequest `json:"intentRequest"`
}

type IntentRequest struct {
	Type     string   `json:"type,omitempty"`
	Selector Selector `json:"selector"`
	Rule     Rule     `json:"rule,omitempty"`
}

type Selector struct {
	Match []Match  `json:"match,omitempty"`
	CEL   []string `json:"cel,omitempty"`
}

type Match struct {
	Kind        string            `json:"kind,omitempty"`
	Condition   string            `json:"condition,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Name        string            `json:"name,omitempty"`
	MatchLabels map[string]string `json:"matchLabels,omitempty"`
}

type Rule struct {
	Action string `json:"action,omitempty"`
	//Mode   string        `json:"mode,omitempty"`

	From []NetPolDetail `json:"from,omitempty"`
	To   []NetPolDetail `json:"to,omitempty"`

	ActionPoint []ActionPoint `json:"actionPoint,omitempty"`
}

type NetPolDetail struct {
	Kind string `json:"kind"`
	// Namespace string `json:"namespace,omitempty"`
	//Endpoint  string            `json:"endpoint,omitempty"`
	Labels   map[string]string `json:"labels,omitempty"`
	Args     []string          `json:"args,omitempty"`
	Port     string            `json:"port,omitempty"`
	Protocol string            `json:"protocol,omitempty"`
}

type ActionPoint struct {
	SubType string `json:"subType"`

	// http
	Headers []EventHeader `json:"headers,omitempty"`

	// cluster
	Precondition []EventFilter `json:"precondition,omitempty"`
	Condition    []EventFilter `json:"conditions,omitempty"`

	Resource EventMatchResource `json:"resource,omitempty"`
}

type EventMatchResource struct {
	// netpol -  http
	// path -> syspol O
	Path    []string `json:"path,omitempty"`
	Methods []string `json:"methods,omitempty"`

	// syspol
	Syscall   string `json:"syscall,omitempty"`
	Subsystem string `json:"subsystem,omitempty"`
	Event     string `json:"event,omitempty"`
	Symbol    string `json:"symbol,omitempty"`

	Dir       string   `json:"dir,omitempty"`
	Pattern   []string `json:"pattern,omitempty"`
	Args      []string `json:"args,omitempty"`
	Protocol  string   `json:"protocol,omitempty"`
	ReadOnly  bool     `json:"readOnly,omitempty"`
	Recursive bool     `json:"recursive,omitempty"`

	// clusterpol
	Kind      string `json:"kind,omitempty"`
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Condition string `json:"condition,omitempty"`
	List      string `json:"list,omitempty"`

	Details []map[string]string `json:"details,omitempty"`
	Filter  []EventFilter       `json:"filter,omitempty"`
	Count   int32               `json:"count,omitempty"`
	Keys    []string            `json:"keys,omitempty"`
	Keyless []Keyless           `json:"keyless,omitempty"`
}

type EventFilter struct {
	Condition string   `json:"Condition,omitempty"`
	Key       string   `json:"key,omitempty"`
	Operator  string   `json:"operator,omitempty"`
	Value     []string `json:"value,omitempty"`
}

type FromSource struct {
	Path      string `json:"path,omitempty"`
	Dir       string `json:"dir,omitempty"`
	Recursive bool   `json:"recursive,omitempty"`
}

type EventHeader struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type Keyless struct {
	Subject string `json:"subject,omitempty"`
	Issuer  string `json:"issuer,omitempty"`
	Url     string `json:"url,omitempty"`
}

// KubeAegisPolicyStatus defines the observed state of KubeAegisPolicy.
type KubeAegisPolicyStatus struct {
	Status            string      `json:"status"`
	LastUpdated       metav1.Time `json:"lastUpdated,omitempty"`
	NumberOfAPs       int32       `json:"numberOfAPs,omitempty"`
	ListofAPs         []string    `json:"listOfAPs,omitempty"`
	NumberOfResources int32       `json:"numberOfResources,omitempty"`
	ListofResources   []string    `json:"listOfResources,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
//+kubebuilder:resource:shortName="kap"
//+kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.status"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
//+kubebuilder:printcolumn:name="Policies",type="string",JSONPath=".status.listOfAPs"
//+kubebuilder:printcolumn:name="Number of APs",type="integer",JSONPath=".status.numberOfAPs"
//+kubebuilder:printcolumn:name="Resources",type="string",JSONPath=".status.listOfResources"
//+kubebuilder:printcolumn:name="Number of Resources",type="integer",JSONPath=".status.numberOfResources"
//+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KubeAegisPolicy is the Schema for the kubeaegispolicies API.
type KubeAegisPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KubeAegisPolicySpec   `json:"spec,omitempty"`
	Status KubeAegisPolicyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KubeAegisPolicyList contains a list of KubeAegisPolicy.
type KubeAegisPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KubeAegisPolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KubeAegisPolicy{}, &KubeAegisPolicyList{})
}
