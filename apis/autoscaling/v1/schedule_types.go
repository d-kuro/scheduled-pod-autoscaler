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
	Suspend bool `json:"suspend"`

	// Description is schedule description.
	// +optional
	Description string `json:"description,omitempty"`

	// TimeZone is the name of the timezone used in the argument of the time.LoadLocation(name string) function.
	// StartTime and EndTime are interpreted as the time in the time zone specified by TimeZone.
	// If not specified, the time will be interpreted as UTC.
	// +optional
	TimeZone string `json:"timeZone,omitempty"`

	// MinReplicas is the lower limit for the number of replicas to which the autoscaler can scale down.
	// It defaults to 1 pod.
	// +kubebuilder:validation:Minimum=1
	// +optional
	MinReplicas *int32 `json:"minReplicas,omitempty"`

	// MaxReplicas is the upper limit for the number of replicas to which the autoscaler can scale up.
	// +kubebuilder:validation:Minimum=1
	// +optional
	MaxReplicas *int32 `json:"maxReplicas,omitempty"`

	// ScheduleType is a type of schedule represented by Weekly,Daily,OneShot.
	// +kubebuiler:validation:Required
	// +kubebuilder:validation:Enum=Monthly;Weekly;Daily;OneShot
	ScheduleType ScheduleType `json:"type"`

	// StartDayOfWeek is scaling start day of week.
	// +kubebuilder:validation:Enum=Monday;Tuesday;Wednesday;Thursday;Friday;Saturday;Sunday
	// +optional
	StartDayOfWeek string `json:"startDayOfWeek"`

	// EndDayOfWeek is scaling end day of week.
	// +kubebuilder:validation:Enum=Monday;Tuesday;Wednesday;Thursday;Friday;Saturday;Sunday
	// +optional
	EndDayOfWeek string `json:"endDayOfWeek"`

	// StartTime is scaling start time. Defined in RFC3339 based format.
	// Different formats are evaluated depending on ScheduleType.
	// e.g. OneShot(yyyy-MM-ddTHH:mm), Monthly(ddTHH:mm), Weekly(HH:mm), Daily(HH:mm)
	// +kubebuiler:validation:Required
	StartTime string `json:"startTime"`

	// EndTime is scaling end time. Defined in RFC3339 based format.
	// Different formats are evaluated depending on ScheduleType.
	// e.g. OneShot(yyyy-MM-ddTHH:mm), Monthly(MM-ddTHH:mm), Weekly(HH:mm), Daily(HH:mm)
	// +kubebuiler:validation:Required
	EndTime string `json:"endTime"`
}

type ScheduleType string

const (
	Monthly ScheduleType = "Monthly"
	Weekly  ScheduleType = "Weekly"
	Daily   ScheduleType = "Daily"
	OneShot ScheduleType = "OneShot"
)

type ScheduleConditionType string

const (
	ScheduleAvailable   ScheduleConditionType = "Available"
	ScheduleSuspend     ScheduleConditionType = "Suspend"
	ScheduleProgressing ScheduleConditionType = "Progressing"
	ScheduleDegraded    ScheduleConditionType = "Degraded"
	ScheduleCompleted   ScheduleConditionType = "Completed"
)

func (s ScheduleSpec) IsCompleted(now time.Time) (bool, error) {
	if s.ScheduleType != OneShot {
		return false, nil
	}

	endTime, err := time.Parse("2006-01-02T15:04", s.EndTime)
	if err != nil {
		return false, err
	}

	return endTime.UTC().After(now.UTC()), nil
}

// ScheduleStatus defines the observed state of Schedule.
type ScheduleStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// LastTransitionTime is the last time the condition transitioned from one status to another.
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Format=date-time
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`

	// Condition is schedule status type.
	// +optional
	Condition ScheduleConditionType `json:"condition,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="REFERENCE",type=string,JSONPath=`.spec.scaleTargetRef.name`,priority=0
// +kubebuilder:printcolumn:name="TYPE",type=string,JSONPath=`.spec.type`,priority=0
// +kubebuilder:printcolumn:name="STARTTIME",type=string,JSONPath=`.spec.startTime`,priority=0
// +kubebuilder:printcolumn:name="ENDTIME",type=string,JSONPath=`.spec.endTime`,priority=0
// +kubebuilder:printcolumn:name="STARTDAYOFWEEK",type=string,JSONPath=`.spec.startDayOfWeek`,priority=0
// +kubebuilder:printcolumn:name="ENDDAYOFWEEK",type=string,JSONPath=`.spec.endDayOfWeek`,priority=0
// +kubebuilder:printcolumn:name="MINPODS",type=integer,JSONPath=`.spec.minReplicas`,priority=1
// +kubebuilder:printcolumn:name="MAXPODS",type=integer,JSONPath=`.spec.maxReplicas`,priority=1
// +kubebuilder:printcolumn:name="STATUS",type=string,JSONPath=`.status.condition`,priority=0
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=`.metadata.creationTimestamp`,priority=0

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
