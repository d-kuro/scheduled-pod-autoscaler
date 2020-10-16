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

package controllers

import (
	"context"

	autoscalingv1 "github.com/d-kuro/scheduled-pod-autoscaler/apis/autoscaling/v1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ScheduleReconciler reconciles a Schedule object.
type ScheduleReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=autoscaling.d-kuro.github.io,resources=schedules,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=autoscaling.d-kuro.github.io,resources=schedules/status,verbs=get;update;patch

func (r *ScheduleReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("schedule", req.NamespacedName)

	var schedule autoscalingv1.Schedule
	if err := r.Get(ctx, req.NamespacedName, &schedule); err != nil {
		log.Error(err, "unable to fetch Schedule")

		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	namespacedName := types.NamespacedName{
		Namespace: schedule.Namespace,
		Name:      schedule.Spec.ScaleTargetRef.Name,
	}

	var spa autoscalingv1.ScheduledPodAutoscaler
	if err := r.Get(ctx, namespacedName, &spa); err != nil {
		log.Error(err, "unable to fetch ScheduledPodAutoscaler", "namespacedName", namespacedName)

		return ctrl.Result{}, err
	}

	if existing := metav1.GetControllerOf(&schedule); existing == nil {
		if err := ctrl.SetControllerReference(&spa, &schedule, r.Scheme); err != nil {
			log.Error(err, "unable to set ownerReference", "schedule", schedule)

			return ctrl.Result{}, err
		}

		if err := r.Update(ctx, &schedule, &client.UpdateOptions{}); err != nil {
			log.Error(err, "unable to update schedule", "schedule", schedule)

			return ctrl.Result{}, err
		}

		log.Info("successfully update Schedule", "schedule", schedule)
	}

	if schedule.Spec.Suspend {
		if updated := setSuspendStatus(&schedule); updated {
			r.Recorder.Event(&schedule, corev1.EventTypeNormal, "Updated", "The schedule was updated.")

			if err := r.Status().Update(ctx, &schedule); err != nil {
				log.Error(err, "unable to update schedule status", "schedule", schedule)

				return ctrl.Result{}, err
			}
		}

		return ctrl.Result{}, nil
	}

	if updated := setAvailableStatus(&schedule); updated {
		r.Recorder.Event(&schedule, corev1.EventTypeNormal, "Updated", "The schedule was updated.")

		if err := r.Status().Update(ctx, &schedule); err != nil {
			log.Error(err, "unable to update schedule status", "schedule", schedule)

			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func setSuspendStatus(schedule *autoscalingv1.Schedule) bool {
	updated := false

	currentSuspendCond := findCondition(schedule.Status.Conditions, string(autoscalingv1.ScheduleSuspend))
	if currentSuspendCond == nil || currentSuspendCond.Status != autoscalingv1.ConditionTrue {
		setCondition(&schedule.Status.Conditions, autoscalingv1.Condition{
			Type:    string(autoscalingv1.ScheduleSuspend),
			Status:  autoscalingv1.ConditionTrue,
			Reason:  "SuspendScheduling",
			Message: "Suspend to scheduled scheduling.",
		})

		schedule.Status.Phase = autoscalingv1.ScheduleSuspend

		updated = true
	}

	currentAvailableCond := findCondition(schedule.Status.Conditions, string(autoscalingv1.ScheduleAvailable))
	if currentAvailableCond == nil || currentAvailableCond.Status != autoscalingv1.ConditionFalse {
		setCondition(&schedule.Status.Conditions, autoscalingv1.Condition{
			Type:    string(autoscalingv1.ScheduleAvailable),
			Status:  autoscalingv1.ConditionFalse,
			Reason:  "SuspendScheduling",
			Message: "Suspend to scheduled scheduling.",
		})

		updated = true
	}

	return updated
}

func setAvailableStatus(schedule *autoscalingv1.Schedule) bool {
	updated := false

	currentAvailableCond := findCondition(schedule.Status.Conditions, string(autoscalingv1.ScheduleAvailable))
	if currentAvailableCond == nil || currentAvailableCond.Status != autoscalingv1.ConditionTrue {
		setCondition(&schedule.Status.Conditions, autoscalingv1.Condition{
			Type:    string(autoscalingv1.ScheduleAvailable),
			Status:  autoscalingv1.ConditionTrue,
			Reason:  "SchedulingAvailable",
			Message: "Available to scheduled scheduling.",
		})

		schedule.Status.Phase = autoscalingv1.ScheduleAvailable

		updated = true
	}

	currentSuspendCond := findCondition(schedule.Status.Conditions, string(autoscalingv1.ScheduleSuspend))
	if currentSuspendCond == nil || currentSuspendCond.Status != autoscalingv1.ConditionFalse {
		setCondition(&schedule.Status.Conditions, autoscalingv1.Condition{
			Type:    string(autoscalingv1.ScheduleSuspend),
			Status:  autoscalingv1.ConditionFalse,
			Reason:  "SchedulingAvailable",
			Message: "Available to scheduled scheduling.",
		})

		updated = true
	}

	return updated
}

func (r *ScheduleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&autoscalingv1.Schedule{}).
		Complete(r)
}
