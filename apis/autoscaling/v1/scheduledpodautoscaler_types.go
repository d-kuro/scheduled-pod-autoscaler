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
	"time"

	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ScheduledPodAutoscalerSpec defines the desired state of ScheduledPodAutoscaler
type ScheduledPodAutoscalerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// HorizontalPodAutoscalerSpec is HorizontalPodAutoscaler v2beta2 API spec.
	// +kubebuilder:validation:Required
	HorizontalPodAutoscalerSpec autoscalingv2beta2.HorizontalPodAutoscalerSpec `json:"horizontalPodAutoscalerSpec"`

	// ScheduleList is list of schedule info.
	// +optional
	ScheduleSpecList []ScheduleSpec `json:"schedule"`
}

// ScheduleSpec is schedule info.
type ScheduleSpec struct {
	// Suspend indicates whether to suspend this schedule.
	// +optional
	Suspend bool `json:"suspend,omitempty"`

	// Name is schedule name.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Description is schedule description.
	// +optional
	Description string `json:"description,omitempty"`

	// MinReplicas is the lower limit for the number of replicas to which the autoscaler can scale down.
	// It defaults to 1 pod.
	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum=1
	// +optional
	MinReplicas *int32 `json:"minReplicas,omitempty"`

	// MaxReplicas is the upper limit for the number of replicas to which the autoscaler can scale up.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=1
	MaxReplicas int32 `json:"maxReplicas"`

	// Metrics contains the specifications for which to use to calculate the desired replica count.
	// +optional
	Metrics []autoscalingv2beta2.MetricSpec `json:"metrics,omitempty"`

	// Behavior configures the scaling behavior of the target in both Up and Down directions.
	// +optional
	// Behavior *autoscalingv2beta2.HorizontalPodAutoscalerBehavior `json:"behavior,omitempty"`

	// StartDayOfWeek is scaling start day of week.
	// +kubebuilder:validation:Enum=Monday;Tuesday;Wednesday;Thursday;Friday;Saturday;Sunday
	// +optional
	StartDayOfWeek string `json:"startDayOfWeek"`

	// EndDayOfWeek is scaling end day of week.
	// +kubebuilder:validation:Enum=Monday;Tuesday;Wednesday;Thursday;Friday;Saturday;Sunday
	// +optional
	EndDayOfWeek string `json:"endDayOfWeek"`

	// StartTime is scaling start time.
	// +kubebuiler:validation:Required
	StartTime string `json:"startTime"`

	// EndTime is scaling end time.
	// +kubebuiler:validation:Required
	EndTime string `json:"endTime"`
}

var Weekdays = map[string]time.Weekday{
	"Monday":    time.Weekday(1),
	"Tuesday":   time.Weekday(2),
	"Wednesday": time.Weekday(3),
	"Thursday":  time.Weekday(4),
	"Friday":    time.Weekday(5),
	"Saturday":  time.Weekday(6),
	"Sunday":    time.Weekday(0),
}

// ScheduledPodAutoscalerStatus defines the observed state of ScheduledPodAutoscaler
type ScheduledPodAutoscalerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// ScheduledPodAutoscaler is the Schema for the scheduledpodautoscalers API
type ScheduledPodAutoscaler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ScheduledPodAutoscalerSpec   `json:"spec,omitempty"`
	Status ScheduledPodAutoscalerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ScheduledPodAutoscalerList contains a list of ScheduledPodAutoscaler
type ScheduledPodAutoscalerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ScheduledPodAutoscaler `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ScheduledPodAutoscaler{}, &ScheduledPodAutoscalerList{})
}
