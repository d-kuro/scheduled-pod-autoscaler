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

package v1

import (
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ScheduledPodAutoscalerSpec defines the desired state of ScheduledPodAutoscaler.
type ScheduledPodAutoscalerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// HorizontalPodAutoscalerSpec is HorizontalPodAutoscaler v2beta2 API spec.
	// ref: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/#horizontalpodautoscaler-v2beta2-autoscaling
	// +kubebuilder:validation:Required
	HorizontalPodAutoscalerSpec autoscalingv2beta2.HorizontalPodAutoscalerSpec `json:"horizontalPodAutoscalerSpec"`
}

type ScheduledPodAutoscalerConditionType string

const (
	ScheduledPodAutoscalerAvailable ScheduledPodAutoscalerConditionType = "Available"
	ScheduledPodAutoscalerDegraded  ScheduledPodAutoscalerConditionType = "Degraded"
)

// ScheduledPodAutoscalerStatus defines the observed state of ScheduledPodAutoscaler.
type ScheduledPodAutoscalerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// LastTransitionTime is the last time the condition transitioned from one status to another.
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Format=date-time
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`

	// Condition is schedule status type.
	// +optional
	Condition ScheduledPodAutoscalerConditionType `json:"condition,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=spa
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="MINPODS",type=integer,JSONPath=`.spec.horizontalPodAutoscalerSpec.minReplicas`,priority=0
// +kubebuilder:printcolumn:name="MAXPODS",type=integer,JSONPath=`.spec.horizontalPodAutoscalerSpec.maxReplicas`,priority=0
// +kubebuilder:printcolumn:name="STATUS",type=string,JSONPath=`.status.condition`,priority=0
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp",priority=0

// ScheduledPodAutoscaler is the Schema for the scheduledpodautoscalers API.
type ScheduledPodAutoscaler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ScheduledPodAutoscalerSpec   `json:"spec,omitempty"`
	Status ScheduledPodAutoscalerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ScheduledPodAutoscalerList contains a list of ScheduledPodAutoscaler.
type ScheduledPodAutoscalerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ScheduledPodAutoscaler `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ScheduledPodAutoscaler{}, &ScheduledPodAutoscalerList{})
}
