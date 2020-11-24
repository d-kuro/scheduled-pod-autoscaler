package controllers

import (
	"context"
	"fmt"
	"time"

	autoscalingv1 "github.com/d-kuro/scheduled-pod-autoscaler/apis/autoscaling/v1"
	"github.com/d-kuro/scheduled-pod-autoscaler/controllers/autoscaling/internal/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	hpav2beta2 "k8s.io/api/autoscaling/v2beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = ginkgo.Describe("ScheduledPodAutoscaler controller", func() {
	ginkgo.Context("when creating ScheduledPodAutoscaler resource", func() {
		ginkgo.It("should create HPA", func() {
			const (
				name = "create-hpa-test"
			)

			ctx := context.Background()
			spa := newScheduledPodAutoscaler(name)

			err := k8sClient.Create(ctx, spa)
			gomega.Expect(err).Should(gomega.Succeed())

			var createdSPA autoscalingv1.ScheduledPodAutoscaler
			var createdHPA hpav2beta2.HorizontalPodAutoscaler

			gomega.Eventually(func() error {
				if err := k8sClient.Get(ctx, client.ObjectKey{Name: name, Namespace: defaultTestNamespace}, &createdHPA); err != nil {
					return err
				}

				if diff := cmp.Diff(spa.Spec.HorizontalPodAutoscalerSpec, createdHPA.Spec); diff != "" {
					return fmt.Errorf("created HPA mismatch (-want +got):\\n%s", diff)
				}

				if err := k8sClient.Get(ctx, client.ObjectKey{Name: name, Namespace: defaultTestNamespace}, &createdSPA); err != nil {
					return err
				}

				if createdSPA.Status.Condition == "" {
					return fmt.Errorf("condition not found")
				}

				if createdSPA.Status.Condition != autoscalingv1.ScheduledPodAutoscalerAvailable {
					return fmt.Errorf("condition not available: %s", createdSPA.Status.Condition)
				}

				return nil
			}, /*timeout*/ defaultTestTimeout /*pollingInterval*/, defaultTestPollingInterval).Should(gomega.Succeed())
		})
		ginkgo.It("should scaling with daily scheduled scaling", func() {
			const (
				name                              = "scheduled-scaling-test"
				scheduledPodAutoscalerMinReplicas = 1
				scheduledPodAutoscalerMaxReplicas = 3
				scheduleMinReplicas               = 5
				scheduleMaxReplicas               = 10
			)

			ctx := context.Background()
			now := time.Now().UTC()
			spa := newScheduledPodAutoscaler(name,
				WithScheduledPodAutoscalerMinReplicas(scheduledPodAutoscalerMinReplicas),
				WithScheduledPodAutoscalerMaxReplicas(scheduledPodAutoscalerMaxReplicas))

			// Target scheduled scaling with start at the current time and end in one hour
			start := now.Format("15:04")
			end := now.Add(time.Hour * 1).Format("15:04")
			schedule := newSchedule(name,
				WithScheduleMinReplicas(scheduleMinReplicas),
				WithScheduleMaxReplicas(scheduleMaxReplicas),
				WithScheduleType(autoscalingv1.Daily),
				WithScheduleStartTime(start),
				WithScheduleEndTime(end))

			err := k8sClient.Create(ctx, spa)
			gomega.Expect(err).Should(gomega.Succeed())

			err = k8sClient.Create(ctx, schedule)
			gomega.Expect(err).Should(gomega.Succeed())

			var createdHPA hpav2beta2.HorizontalPodAutoscaler
			var createdSchedule autoscalingv1.Schedule
			gomega.Eventually(func() error {
				if err := k8sClient.Get(ctx, client.ObjectKey{Name: name, Namespace: defaultTestNamespace}, &createdHPA); err != nil {
					return err
				}

				if createdHPA.Spec.MinReplicas == nil {
					return fmt.Errorf("created HPA minReplicas mismatch: want: %d, got: nil", scheduleMinReplicas)
				}

				if *createdHPA.Spec.MinReplicas != int32(scheduleMinReplicas) {
					return fmt.Errorf("created HPA minReplicas mismatch: want: %d, got: %d",
						scheduleMinReplicas, *createdHPA.Spec.MinReplicas)
				}

				if createdHPA.Spec.MaxReplicas != int32(scheduleMaxReplicas) {
					return fmt.Errorf("created HPA maxReplicas mismatch: want: %d, got: %d",
						scheduleMaxReplicas, createdHPA.Spec.MaxReplicas)
				}

				if err := k8sClient.Get(ctx, client.ObjectKey{Name: name, Namespace: defaultTestNamespace}, &createdSchedule); err != nil {
					return err
				}

				if createdSchedule.Status.Condition != autoscalingv1.ScheduleProgressing {
					return fmt.Errorf("schedule condition mismatch: want: %s, got: %s",
						autoscalingv1.ScheduleProgressing, createdSchedule.Status.Condition)
				}

				return nil
			}, /*timeout*/ defaultTestTimeout /*pollingInterval*/, defaultTestPollingInterval).Should(gomega.Succeed())
		})
		ginkgo.It("should suspend scheduled scaling", func() {
			const (
				name                              = "scheduled-scaling-suspend-test"
				namespace                         = "default"
				scheduledPodAutoscalerMinReplicas = 1
				scheduledPodAutoscalerMaxReplicas = 3
				scheduleMinReplicas               = 5
				scheduleMaxReplicas               = 10
			)

			ctx := context.Background()
			now := time.Now().UTC()
			spa := newScheduledPodAutoscaler(name,
				WithScheduledPodAutoscalerMinReplicas(scheduledPodAutoscalerMinReplicas),
				WithScheduledPodAutoscalerMaxReplicas(scheduledPodAutoscalerMaxReplicas))

			// Target scheduled scaling with start at the current time and end in one hour
			start := now.Format("15:04")
			end := now.Add(time.Hour * 1).Format("15:04")
			schedule := newSchedule(name,
				WithScheduleSuspend(true), // suspend for scheduled scaling
				WithScheduleMinReplicas(scheduleMinReplicas),
				WithScheduleMaxReplicas(scheduleMaxReplicas),
				WithScheduleType(autoscalingv1.Daily),
				WithScheduleStartTime(start),
				WithScheduleEndTime(end))

			err := k8sClient.Create(ctx, spa)
			gomega.Expect(err).Should(gomega.Succeed())

			err = k8sClient.Create(ctx, schedule)
			gomega.Expect(err).Should(gomega.Succeed())

			var createdHPA hpav2beta2.HorizontalPodAutoscaler
			var createdSchedule autoscalingv1.Schedule
			gomega.Eventually(func() error {
				if err := k8sClient.Get(ctx, client.ObjectKey{Name: name, Namespace: namespace}, &createdHPA); err != nil {
					return err
				}

				if createdHPA.Spec.MinReplicas == nil {
					return fmt.Errorf("created HPA minReplicas mismatch: want: %d, got: nil", scheduleMinReplicas)
				}

				if *createdHPA.Spec.MinReplicas != int32(scheduledPodAutoscalerMinReplicas) {
					return fmt.Errorf("created HPA minReplicas mismatch: want: %d, got: %d",
						scheduledPodAutoscalerMinReplicas, *createdHPA.Spec.MinReplicas)
				}

				if createdHPA.Spec.MaxReplicas != int32(scheduledPodAutoscalerMaxReplicas) {
					return fmt.Errorf("created HPA maxReplicas mismatch: want: %d, got: %d",
						scheduledPodAutoscalerMaxReplicas, createdHPA.Spec.MaxReplicas)
				}

				if err := k8sClient.Get(ctx, client.ObjectKey{Name: name, Namespace: namespace}, &createdSchedule); err != nil {
					return err
				}

				if createdSchedule.Status.Condition != autoscalingv1.ScheduleSuspend {
					return fmt.Errorf("schedule condition mismatch: want: %s, got: %s",
						autoscalingv1.ScheduleSuspend, createdSchedule.Status.Condition)
				}

				return nil
			}, /*timeout*/ defaultTestTimeout /*pollingInterval*/, defaultTestPollingInterval).Should(gomega.Succeed())
		})
		ginkgo.It("should completed scheduled scaling", func() {
			const (
				name                              = "scheduled-scaling-completed-test"
				namespace                         = "default"
				scheduledPodAutoscalerMinReplicas = 1
				scheduledPodAutoscalerMaxReplicas = 3
				scheduleMinReplicas               = 5
				scheduleMaxReplicas               = 10
			)

			ctx := context.Background()
			now := time.Now().UTC()
			spa := newScheduledPodAutoscaler(name,
				WithScheduledPodAutoscalerMinReplicas(scheduledPodAutoscalerMinReplicas),
				WithScheduledPodAutoscalerMaxReplicas(scheduledPodAutoscalerMaxReplicas))

			// A one-shot schedule with scaling completed
			start := now.AddDate(0, 0, -1).Format("2006-01-02T15:04")
			end := now.AddDate(0, 0, -10).Format("2006-01-02T15:04")
			schedule := newSchedule(name,
				WithScheduleType(autoscalingv1.OneShot),
				WithScheduleStartTime(start),
				WithScheduleEndTime(end))

			err := k8sClient.Create(ctx, spa)
			gomega.Expect(err).Should(gomega.Succeed())

			err = k8sClient.Create(ctx, schedule)
			gomega.Expect(err).Should(gomega.Succeed())

			var createdHPA hpav2beta2.HorizontalPodAutoscaler
			var createdSchedule autoscalingv1.Schedule
			gomega.Eventually(func() error {
				if err := k8sClient.Get(ctx, client.ObjectKey{Name: name, Namespace: namespace}, &createdHPA); err != nil {
					return err
				}

				if createdHPA.Spec.MinReplicas == nil {
					return fmt.Errorf("created HPA minReplicas mismatch: want: %d, got: nil", scheduleMinReplicas)
				}

				if *createdHPA.Spec.MinReplicas != int32(scheduledPodAutoscalerMinReplicas) {
					return fmt.Errorf("created HPA minReplicas mismatch: want: %d, got: %d",
						scheduledPodAutoscalerMinReplicas, *createdHPA.Spec.MinReplicas)
				}

				if createdHPA.Spec.MaxReplicas != int32(scheduledPodAutoscalerMaxReplicas) {
					return fmt.Errorf("created HPA maxReplicas mismatch: want: %d, got: %d",
						scheduledPodAutoscalerMaxReplicas, createdHPA.Spec.MaxReplicas)
				}

				if err := k8sClient.Get(ctx, client.ObjectKey{Name: name, Namespace: namespace}, &createdSchedule); err != nil {
					return err
				}

				if createdSchedule.Status.Condition != autoscalingv1.ScheduleCompleted {
					return fmt.Errorf("schedule condition mismatch: want: %s, got: %s",
						autoscalingv1.ScheduleCompleted, createdSchedule.Status.Condition)
				}

				return nil
			}, /*timeout*/ defaultTestTimeout /*pollingInterval*/, defaultTestPollingInterval).Should(gomega.Succeed())
		})
	})
})

const (
	defaultSPAMinReplicas = 1
	defaultSPAMaxReplicas = 3
)

func newScheduledPodAutoscaler(name string, options ...func(*autoscalingv1.ScheduledPodAutoscaler)) *autoscalingv1.ScheduledPodAutoscaler {
	spa := &autoscalingv1.ScheduledPodAutoscaler{
		TypeMeta: metav1.TypeMeta{
			APIVersion: autoscalingv1.GroupVersion.String(),
			Kind:       "ScheduledPodAutoscaler",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: defaultTestNamespace,
		},
		Spec: autoscalingv1.ScheduledPodAutoscalerSpec{
			HorizontalPodAutoscalerSpec: hpav2beta2.HorizontalPodAutoscalerSpec{
				ScaleTargetRef: hpav2beta2.CrossVersionObjectReference{
					APIVersion: "apps/v1",
					Kind:       "Deployment",
					Name:       name,
				},
				MinReplicas: testutil.ToPointerInt32(defaultSPAMinReplicas),
				MaxReplicas: defaultSPAMaxReplicas,
				Metrics: []hpav2beta2.MetricSpec{
					{
						Type: "Resource",
						Resource: &hpav2beta2.ResourceMetricSource{
							Name: "cpu",
							Target: hpav2beta2.MetricTarget{
								Type:               "Utilization",
								AverageUtilization: testutil.ToPointerInt32(50),
							},
						},
					},
				},
			},
		},
	}

	for _, option := range options {
		option(spa)
	}

	return spa
}

func WithScheduledPodAutoscalerMinReplicas(value int) func(*autoscalingv1.ScheduledPodAutoscaler) {
	return func(spa *autoscalingv1.ScheduledPodAutoscaler) {
		spa.Spec.HorizontalPodAutoscalerSpec.MinReplicas = testutil.ToPointerInt32(value)
	}
}

func WithScheduledPodAutoscalerMaxReplicas(value int) func(*autoscalingv1.ScheduledPodAutoscaler) {
	return func(spa *autoscalingv1.ScheduledPodAutoscaler) {
		spa.Spec.HorizontalPodAutoscalerSpec.MaxReplicas = int32(value)
	}
}
