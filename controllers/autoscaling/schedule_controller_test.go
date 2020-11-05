package controllers

import (
	"context"
	"fmt"
	"time"

	autoscalingv1 "github.com/d-kuro/scheduled-pod-autoscaler/apis/autoscaling/v1"
	"github.com/d-kuro/scheduled-pod-autoscaler/controllers/autoscaling/internal/testutil"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	hpav2beta2 "k8s.io/api/autoscaling/v2beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = ginkgo.Describe("Schedule controller", func() {
	ginkgo.Context("when creating Schedule resource", func() {
		ginkgo.It("should set ownerReference", func() {
			const (
				name      = "schedule-controller-test"
				namespace = "default"
			)

			ctx := context.Background()
			now := time.Now().UTC()
			spa := newScheduledPodAutoscaler(name, namespace)

			// Set a future time and prevent it from being scheduled scaling
			start := now.AddDate(0, 0, 1).Format("2006-01-02T15:04")
			end := now.AddDate(0, 0, 10).Format("2006-01-02T15:04")
			schedule := newSchedule(name, namespace,
				WithScheduleType(autoscalingv1.OneShot),
				WithScheduleStartTime(start),
				WithScheduleEndTime(end))

			err := k8sClient.Create(ctx, spa)
			gomega.Expect(err).Should(gomega.Succeed())

			err = k8sClient.Create(ctx, schedule)
			gomega.Expect(err).Should(gomega.Succeed())

			var createdSchedule autoscalingv1.Schedule
			gomega.Eventually(func() error {
				if err := k8sClient.Get(ctx, client.ObjectKey{Name: name, Namespace: namespace}, &createdSchedule); err != nil {
					return err
				}

				if createdSchedule.Status.Condition != autoscalingv1.ScheduleAvailable {
					return fmt.Errorf("schedule condition mismatch: want: %s, got: %s",
						autoscalingv1.ScheduleAvailable, createdSchedule.Status.Condition)
				}

				if existing := metav1.GetControllerOf(&createdSchedule); existing == nil {
					return fmt.Errorf("ownerReference not found")
				}

				return nil
			}, /*timeout*/ time.Second*1 /*pollingInterval*/, time.Millisecond*100).Should(gomega.Succeed())
		})
	})
})

const (
	defaultScheduleMinReplicas = 3
	defaultScheduleMaxReplicas = 10
	defaultScheduleType        = "Daily"
	defaultScheduleStartTime   = "00:00"
	defaultScheduleEndTime     = "12:00"
)

func newSchedule(name string, namespace string, options ...func(*autoscalingv1.Schedule)) *autoscalingv1.Schedule {
	schedule := &autoscalingv1.Schedule{
		TypeMeta: metav1.TypeMeta{
			APIVersion: autoscalingv1.GroupVersion.String(),
			Kind:       "Schedule",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: autoscalingv1.ScheduleSpec{
			ScaleTargetRef: hpav2beta2.CrossVersionObjectReference{
				APIVersion: "autoscaling.d-kuro.github.io/v1",
				Kind:       "ScheduledPodAutoscaler",
				Name:       name,
			},
			ScheduleType: defaultScheduleType,
			MinReplicas:  testutil.ToPointerInt32(defaultScheduleMinReplicas),
			MaxReplicas:  testutil.ToPointerInt32(defaultScheduleMaxReplicas),
			StartTime:    defaultScheduleStartTime,
			EndTime:      defaultScheduleEndTime,
		},
	}

	for _, option := range options {
		option(schedule)
	}

	return schedule
}

func WithScheduleType(value autoscalingv1.ScheduleType) func(*autoscalingv1.Schedule) {
	return func(schedule *autoscalingv1.Schedule) {
		schedule.Spec.ScheduleType = value
	}
}

func WithScheduleStartTime(t string) func(*autoscalingv1.Schedule) {
	return func(schedule *autoscalingv1.Schedule) {
		schedule.Spec.StartTime = t
	}
}

func WithScheduleEndTime(t string) func(*autoscalingv1.Schedule) {
	return func(schedule *autoscalingv1.Schedule) {
		schedule.Spec.EndTime = t
	}
}

func WithScheduleMinReplicas(value int) func(*autoscalingv1.Schedule) {
	return func(schedule *autoscalingv1.Schedule) {
		schedule.Spec.MinReplicas = testutil.ToPointerInt32(value)
	}
}

func WithScheduleMaxReplicas(value int) func(*autoscalingv1.Schedule) {
	return func(schedule *autoscalingv1.Schedule) {
		schedule.Spec.MaxReplicas = testutil.ToPointerInt32(value)
	}
}
