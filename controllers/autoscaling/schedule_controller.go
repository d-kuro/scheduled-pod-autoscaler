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
	apierrors "k8s.io/apimachinery/pkg/api/errors"
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

func (r *ScheduleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("schedule", req.NamespacedName)

	var schedule autoscalingv1.Schedule
	if err := r.Get(ctx, req.NamespacedName, &schedule); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}

		log.Error(err, "unable to fetch Schedule")

		return ctrl.Result{}, err
	}

	if schedule.DeletionTimestamp != nil {
		return ctrl.Result{}, nil
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
		if schedule.Status.Condition != autoscalingv1.ScheduleSuspend {
			if err := r.updateScheduleStatus(ctx, log, schedule, autoscalingv1.ScheduleSuspend); err != nil {
				return ctrl.Result{}, err
			}
		}

		return ctrl.Result{}, nil
	}

	if schedule.Status.Condition != autoscalingv1.ScheduleSuspend &&
		schedule.Status.Condition != autoscalingv1.ScheduleCompleted {
		if err := r.updateScheduleStatus(ctx, log, schedule, autoscalingv1.ScheduleAvailable); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *ScheduleReconciler) updateScheduleStatus(ctx context.Context, log logr.Logger,
	schedule autoscalingv1.Schedule, newCondition autoscalingv1.ScheduleConditionType) error {
	if updated := setScheduleCondition(&schedule.Status, newCondition); updated {
		r.Recorder.Event(&schedule, corev1.EventTypeNormal, "Updated", "The schedule was updated.")

		if err := r.Status().Update(ctx, &schedule); err != nil {
			log.Error(err, "unable to update schedule status", "schedule", schedule)

			return err
		}
	}

	return nil
}

func setScheduleCondition(status *autoscalingv1.ScheduleStatus, newCondition autoscalingv1.ScheduleConditionType) bool {
	updated := false

	if status.Condition == newCondition {
		return updated
	}

	status.Condition = newCondition
	status.LastTransitionTime = metav1.Now()
	updated = true

	return updated
}

func (r *ScheduleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&autoscalingv1.Schedule{}).
		Complete(r)
}
