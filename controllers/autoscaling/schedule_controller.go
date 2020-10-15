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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ScheduleReconciler reconciles a Schedule object.
type ScheduleReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=autoscaling.d-kuro.github.io,resources=schedules,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=autoscaling.d-kuro.github.io,resources=schedules/status,verbs=get;update;patch

func (r *ScheduleReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("schedule", req.NamespacedName)

	var schedule autoscalingv1.Schedule
	if err := r.Get(ctx, req.NamespacedName, &schedule); err != nil {
		log.Error(err, "unable to fetch Schedule")

		return ctrl.Result{}, err
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

	if err := ctrl.SetControllerReference(&spa, &schedule, r.Scheme); err != nil {
		log.Error(err, "unable to set ownerReference", "schedule", schedule)

		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *ScheduleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&autoscalingv1.Schedule{}).
		Complete(r)
}
