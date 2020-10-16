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

// ScheduleSpec defines the desired state of Schedule.
type ScheduleSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// scaleTargetRef points to the target resource to scale, and is used to the pods for which metrics
	// should be collected, as well as to actually change the replica count.
	// +kubebuiler:validation:Required
	ScaleTargetRef autoscalingv2beta2.CrossVersionObjectReference `json:"scaleTargetRef"`

	// Suspend indicates whether to suspend this schedule.
	// +optional
	Suspend bool `json:"suspend,omitempty"`

	// Description is schedule description.
	// +optional
	Description string `json:"description,omitempty"`

	// TimeZoneName is the name of the timezone used in the argument of the time.LoadLocation(name string) function.
	// StartTime and EndTime are interpreted as the time in the time zone specified by TimeZoneName.
	// If not specified, the time will be interpreted as UTC.
	// +optional
	TimeZoneName string `json:"timeZoneName,omitempty"`

	// MinReplicas is the lower limit for the number of replicas to which the autoscaler can scale down.
	// It defaults to 1 pod.
	// +kubebuilder:validation:Minimum=1
	// +optional
	MinReplicas *int32 `json:"minReplicas,omitempty"`

	// MaxReplicas is the upper limit for the number of replicas to which the autoscaler can scale up.
	// +kubebuilder:validation:Minimum=1
	// +optional
	MaxReplicas *int32 `json:"maxReplicas,omitempty"`

	// Metrics contains the specifications for which to use to calculate the desired replica count.
	// +optional
	Metrics []autoscalingv2beta2.MetricSpec `json:"metrics,omitempty"`

	// Behavior configures the scaling behavior of the target in both Up and Down directions.
	// +optional
	Behavior *autoscalingv2beta2.HorizontalPodAutoscalerBehavior `json:"behavior,omitempty"`

	// StartDayOfWeek is scaling start day of week.
	// +kubebuiler:validation:Required
	// +kubebuilder:validation:Enum=Monday;Tuesday;Wednesday;Thursday;Friday;Saturday;Sunday
	StartDayOfWeek string `json:"startDayOfWeek"`

	// EndDayOfWeek is scaling end day of week.
	// +kubebuiler:validation:Required
	// +kubebuilder:validation:Enum=Monday;Tuesday;Wednesday;Thursday;Friday;Saturday;Sunday
	EndDayOfWeek string `json:"endDayOfWeek"`

	// StartTime is scaling start time.
	// +kubebuiler:validation:Required
	StartTime string `json:"startTime"`

	// EndTime is scaling end time.
	// +kubebuiler:validation:Required
	EndTime string `json:"endTime"`
}

// ScheduleStatus defines the observed state of Schedule.
type ScheduleStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="REFERENCE",type=string,JSONPath=`.spec.scaleTargetRef.name`
// +kubebuilder:printcolumn:name="MINPODS",type=integer,JSONPath=`.spec.minReplicas`
// +kubebuilder:printcolumn:name="MAXPODS",type=integer,JSONPath=`.spec.maxReplicas`
// +kubebuilder:printcolumn:name="STARTTIME",type=string,JSONPath=`.spec.startTime`
// +kubebuilder:printcolumn:name="STARTDAYOFWEEK",type=string,JSONPath=`.spec.startDayOfWeek`
// +kubebuilder:printcolumn:name="ENDTIME",type=string,JSONPath=`.spec.endTime`
// +kubebuilder:printcolumn:name="ENDDAYOFWEEK",type=string,JSONPath=`.spec.endDayOfWeek`
// +kubebuilder:printcolumn:name="SUSPEND",type=string,JSONPath=`.spec.suspend`

// Schedule is the Schema for the schedules API.
type Schedule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ScheduleSpec   `json:"spec,omitempty"`
	Status ScheduleStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ScheduleList contains a list of Schedule.
type ScheduleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Schedule `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Schedule{}, &ScheduleList{})
}
