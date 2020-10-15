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
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ScheduledPodAutoscalerReconciler reconciles a ScheduledPodAutoscaler object.
type ScheduledPodAutoscalerReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=autoscaling.d-kuro.github.io,resources=scheduledpodautoscalers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=autoscaling.d-kuro.github.io,resources=scheduledpodautoscalers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=autoscaling,resources=horizontalpodautoscalers,verbs=get;list;watch;create;update;patch;delete

func (r *ScheduledPodAutoscalerReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("scheduledpodautoscaler", req.NamespacedName)

	var scheduledPodAutoscaler autoscalingv1.ScheduledPodAutoscaler
	if err := r.Get(ctx, req.NamespacedName, &scheduledPodAutoscaler); err != nil {
		log.Error(err, "unable to fetch ScheduledPodAutoscaler")

		return ctrl.Result{}, err
	}

	var hpa hpav2beta2.HorizontalPodAutoscaler
	if err := r.Get(ctx, req.NamespacedName, &hpa); apierrors.IsNotFound(err) {
		hpa = hpav2beta2.HorizontalPodAutoscaler{
			ObjectMeta: metav1.ObjectMeta{
				Name:      req.Name,
				Namespace: req.Namespace,
			},
			Spec: scheduledPodAutoscaler.Spec.HorizontalPodAutoscalerSpec,
		}

		if err := ctrl.SetControllerReference(&scheduledPodAutoscaler, &hpa, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}

		if err := r.Create(ctx, &hpa, &client.CreateOptions{}); err != nil {
			return ctrl.Result{}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	}

	if len(scheduledPodAutoscaler.Spec.ScheduleSpecList) >= 2 {
		sort.SliceStable(scheduledPodAutoscaler.Spec.ScheduleSpecList, func(i, j int) bool {
			return scheduledPodAutoscaler.Spec.ScheduleSpecList[i].Name < scheduledPodAutoscaler.Spec.ScheduleSpecList[j].Name
		})
	}

	now := time.Now()

	for _, schedule := range scheduledPodAutoscaler.Spec.ScheduleSpecList {
		isContains, err := schedule.Contains(now)
		if err != nil {
			return ctrl.Result{}, err
		}

		if isContains {
			hpa.Spec.MaxReplicas = schedule.MaxReplicas
			hpa.Spec.MinReplicas = schedule.MinReplicas
			hpa.Spec.Metrics = schedule.Metrics

			if err := r.Update(ctx, &hpa, &client.UpdateOptions{}); err != nil {
				return ctrl.Result{}, err
			}

			return ctrl.Result{}, nil
		}
	}

	hpa.Spec = scheduledPodAutoscaler.Spec.HorizontalPodAutoscalerSpec
	if err := r.Update(ctx, &hpa, &client.UpdateOptions{}); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *ScheduledPodAutoscalerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&autoscalingv1.ScheduledPodAutoscaler{}).
		Complete(r)
}
