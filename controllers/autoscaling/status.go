package controllers

import (
	"time"

	autoscalingv1 "github.com/d-kuro/scheduled-pod-autoscaler/apis/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func setCondition(conditions *[]autoscalingv1.Condition, newCondition autoscalingv1.Condition) {
	if conditions == nil {
		conditions = &[]autoscalingv1.Condition{}
	}

	current := findCondition(*conditions, newCondition.Type)
	if current == nil {
		newCondition.LastTransitionTime = metav1.NewTime(time.Now())
		*conditions = append(*conditions, newCondition)

		return
	}

	if current.Status != newCondition.Status {
		current.Status = newCondition.Status
		current.LastTransitionTime = metav1.NewTime(time.Now())
	}

	current.Reason = newCondition.Reason
	current.Message = newCondition.Message
}

func findCondition(conditions []autoscalingv1.Condition, conditionType string) *autoscalingv1.Condition {
	for _, c := range conditions {
		if c.Type == conditionType {
			return &c
		}
	}

	return nil
}
