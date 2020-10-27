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
	"sort"
	"time"

	autoscalingv1 "github.com/d-kuro/scheduled-pod-autoscaler/apis/autoscaling/v1"
	"github.com/go-logr/logr"
	hpav2beta2 "k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ScheduledPodAutoscalerReconciler reconciles a ScheduledPodAutoscaler object.
type ScheduledPodAutoscalerReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=autoscaling.d-kuro.github.io,resources=scheduledpodautoscalers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=autoscaling.d-kuro.github.io,resources=scheduledpodautoscalers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=autoscaling,resources=horizontalpodautoscalers,verbs=get;list;watch;create;update;patch;delete

func (r *ScheduledPodAutoscalerReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("scheduledpodautoscaler", req.NamespacedName)

	var spa autoscalingv1.ScheduledPodAutoscaler
	if err := r.Get(ctx, req.NamespacedName, &spa); err != nil {
		log.Error(err, "unable to fetch ScheduledPodAutoscaler")

		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var hpa hpav2beta2.HorizontalPodAutoscaler
	if err := r.Get(ctx, req.NamespacedName, &hpa); apierrors.IsNotFound(err) {
		log.Info("unable to fetch hpa, try to create one", "namespacedName", req.NamespacedName)

		if err := r.createHPA(ctx, log, spa); err != nil {
			return ctrl.Result{}, err
		}

		log.Info("successfully create HPA", "hpa", hpa)
	} else if err != nil {
		log.Error(err, "unable to fetch HPA", "namespacedName", req.NamespacedName)

		return ctrl.Result{}, err
	}

	updated, err := r.reconcileSchedule(ctx, log, spa, hpa)
	if err != nil {
		log.Error(err, "unable to reconcile")

		return ctrl.Result{}, err
	}

	if !updated {
		hpa.Spec = spa.Spec.HorizontalPodAutoscalerSpec
		if err := r.Update(ctx, &hpa, &client.UpdateOptions{}); err != nil {
			log.Error(err, "unable to update HPA", "hpa", hpa)

			return ctrl.Result{}, err
		}

		if err := r.updateScheduledPodAutoscalerStatus(ctx, log, spa, autoscalingv1.ScheduledPodAutoscalerAvailable); err != nil {
			log.Error(err, "unable to update ScheduledPodAutoscaler status", "scheduledPodAutoscaler", spa)
		}

		log.Info("successfully update HPA", "hpa", hpa)
	}

	return ctrl.Result{Requeue: true, RequeueAfter: 30 * time.Second}, nil
}

func (r *ScheduledPodAutoscalerReconciler) reconcileSchedule(ctx context.Context, log logr.Logger,
	spa autoscalingv1.ScheduledPodAutoscaler, hpa hpav2beta2.HorizontalPodAutoscaler) (bool, error) {
	var schedules autoscalingv1.ScheduleList
	if err := r.List(ctx, &schedules, client.MatchingFields(map[string]string{ownerControllerField: spa.Name})); err != nil {
		log.Error(err, "unable to list child Schedules")

		return false, err
	}

	sortSchedules(schedules.Items, log)

	now := time.Now()
	updated := false

	for _, schedule := range schedules.Items {
		schedule := schedule

		if schedule.Spec.Suspend {
			continue
		}

		isContains, err := schedule.Spec.Contains(now)
		if err != nil {
			log.Error(err, "unable to check contains Schedule")

			return updated, err
		}

		if isContains {
			updated, err = r.updateHPA(ctx, log, schedule, hpa)
			if err != nil {
				return updated, err
			}

			continue
		}

		if err = r.updateScheduleStatus(ctx, log, schedule, autoscalingv1.ScheduleAvailable); err != nil {
			log.Error(err, "unable to update schedule status", "schedule", schedule)
		}
	}

	return updated, nil
}

func (r *ScheduledPodAutoscalerReconciler) createHPA(ctx context.Context, log logr.Logger,
	spa autoscalingv1.ScheduledPodAutoscaler) error {
	hpa := hpav2beta2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      spa.Name,
			Namespace: spa.Namespace,
		},
		Spec: spa.Spec.HorizontalPodAutoscalerSpec,
	}

	if err := ctrl.SetControllerReference(&spa, &hpa, r.Scheme); err != nil {
		log.Error(err, "unable to set ownerReference", "hpa", hpa)

		return err
	}

	if err := r.Create(ctx, &hpa, &client.CreateOptions{}); err != nil {
		log.Info("unable to HPA", "hpa", hpa)

		if err := r.updateScheduledPodAutoscalerStatus(ctx, log, spa, autoscalingv1.ScheduledPodAutoscalerDegraded); err != nil {
			log.Error(err, "unable to update ScheduledPodAutoscaler status", "scheduledPodAutoscaler", spa)
		}

		return err
	}

	return nil
}

func (r *ScheduledPodAutoscalerReconciler) updateHPA(ctx context.Context, log logr.Logger,
	schedule autoscalingv1.Schedule, hpa hpav2beta2.HorizontalPodAutoscaler) (bool, error) {
	updated := false

	if schedule.Spec.MaxReplicas != nil {
		hpa.Spec.MaxReplicas = *schedule.Spec.MaxReplicas
	}

	if schedule.Spec.MinReplicas != nil {
		hpa.Spec.MinReplicas = schedule.Spec.MinReplicas
	}

	if schedule.Spec.Metrics != nil {
		hpa.Spec.Metrics = schedule.Spec.Metrics
	}

	if err := r.Update(ctx, &hpa, &client.UpdateOptions{}); err != nil {
		log.Error(err, "unable to update HPA", "hpa", hpa)

		if err = r.updateScheduleStatus(ctx, log, schedule, autoscalingv1.ScheduleDegraded); err != nil {
			log.Error(err, "unable to update schedule status", "schedule", schedule)
		}

		return updated, err
	}

	updated = true
	log.Info("successfully update HPA", "hpa", hpa)

	if err := r.updateScheduleStatus(ctx, log, schedule, autoscalingv1.ScheduleProgressing); err != nil {
		log.Error(err, "unable to update schedule status", "schedule", schedule)
	}

	return updated, nil
}

func (r *ScheduledPodAutoscalerReconciler) updateScheduleStatus(ctx context.Context, log logr.Logger,
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

func (r *ScheduledPodAutoscalerReconciler) updateScheduledPodAutoscalerStatus(ctx context.Context, log logr.Logger,
	spa autoscalingv1.ScheduledPodAutoscaler, newCondition autoscalingv1.ScheduledPodAutoscalerConditionType) error {
	if updated := setScheduledPodAutoscalerCondition(&spa.Status, newCondition); updated {
		r.Recorder.Event(&spa, corev1.EventTypeNormal, "Updated", "The schedule was updated.")

		if err := r.Status().Update(ctx, &spa); err != nil {
			log.Error(err, "unable to update ScheduledPodAutoscaler status",
				"scheduledPodAutoscaler", spa)

			return err
		}
	}

	return nil
}

func setScheduledPodAutoscalerCondition(
	status *autoscalingv1.ScheduledPodAutoscalerStatus,
	newCondition autoscalingv1.ScheduledPodAutoscalerConditionType) bool {
	updated := false

	if status.Condition == newCondition {
		return updated
	}

	status.Condition = newCondition
	status.LastTransitionTime = metav1.Now()
	updated = true

	return updated
}

// sortSchedules sorts the schedule slice.
// Sort by start day of week.
// If the start day of week are the same. sort by start time.
func sortSchedules(schedules []autoscalingv1.Schedule, log logr.Logger) {
	sort.SliceStable(schedules, func(i, j int) bool {
		if schedules[i].Spec.StartDayOfWeek == schedules[j].Spec.StartDayOfWeek {
			startTime1, err := time.Parse("15:04", schedules[i].Spec.StartTime)
			if err != nil {
				log.Error(err, "unable to parse start time", "schedule", schedules[i])

				return false
			}

			startTime2, err := time.Parse("15:04", schedules[j].Spec.StartTime)
			if err != nil {
				log.Error(err, "unable to parse start time", "schedule", schedules[j])

				return false
			}

			return startTime1.Unix() < startTime2.Unix()
		}

		return schedules[i].Spec.StartDayOfWeek < schedules[j].Spec.StartDayOfWeek
	})
}

const ownerControllerField = ".metadata.controller"

func indexByOwnerScheduledPodAutoscaler(obj runtime.Object) []string {
	schedule := obj.(*autoscalingv1.Schedule)

	owner := metav1.GetControllerOf(schedule)
	if owner == nil {
		return nil
	}

	if owner.APIVersion != autoscalingv1.GroupVersion.String() || owner.Kind != "ScheduledPodAutoscaler" {
		return nil
	}

	return []string{owner.Name}
}

func (r *ScheduledPodAutoscalerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	ctx := context.Background()

	err := mgr.GetFieldIndexer().
		IndexField(ctx, &autoscalingv1.Schedule{}, ownerControllerField, indexByOwnerScheduledPodAutoscaler)
	if err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&autoscalingv1.ScheduledPodAutoscaler{}).
		Owns(&autoscalingv1.Schedule{}).
		Owns(&hpav2beta2.HorizontalPodAutoscaler{}).
		Complete(r)
}
