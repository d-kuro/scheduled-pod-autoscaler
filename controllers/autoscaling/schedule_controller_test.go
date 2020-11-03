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
	const (
		scheduleName = "test"
		namespace    = "default"
	)

	ginkgo.Context("when creating Schedule resource", func() {
		ginkgo.It("should set ownerReference", func() {
			ctx := context.Background()

			spa := &autoscalingv1.ScheduledPodAutoscaler{
				TypeMeta: metav1.TypeMeta{
					APIVersion: autoscalingv1.GroupVersion.String(),
					Kind:       "ScheduledPodAutoscaler",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      scheduleName,
					Namespace: namespace,
				},
				Spec: autoscalingv1.ScheduledPodAutoscalerSpec{
					HorizontalPodAutoscalerSpec: hpav2beta2.HorizontalPodAutoscalerSpec{
						ScaleTargetRef: hpav2beta2.CrossVersionObjectReference{
							APIVersion: "apps/v1",
							Kind:       "Deployment",
							Name:       scheduleName,
						},
						MinReplicas: testutil.ToPointerInt32(1),
						MaxReplicas: 3,
					},
				},
			}

			schedule := &autoscalingv1.Schedule{
				TypeMeta: metav1.TypeMeta{
					APIVersion: autoscalingv1.GroupVersion.String(),
					Kind:       "Schedule",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      scheduleName,
					Namespace: namespace,
				},
				Spec: autoscalingv1.ScheduleSpec{
					ScaleTargetRef: hpav2beta2.CrossVersionObjectReference{
						APIVersion: "autoscaling.d-kuro.github.io/v1",
						Kind:       "ScheduledPodAutoscaler",
						Name:       scheduleName,
					},
					ScheduleType: "Daily",
					MinReplicas:  testutil.ToPointerInt32(5),
					MaxReplicas:  testutil.ToPointerInt32(5),
					StartTime:    "00:00",
					EndTime:      "12:00",
				},
			}

			err := k8sClient.Create(ctx, spa)
			gomega.Expect(err).Should(gomega.Succeed())

			err = k8sClient.Create(ctx, schedule)
			gomega.Expect(err).Should(gomega.Succeed())

			var createdSchedule autoscalingv1.Schedule
			gomega.Eventually(func() error {
				if err := k8sClient.Get(ctx, client.ObjectKey{Name: scheduleName, Namespace: namespace}, &createdSchedule); err != nil {
					return err
				}

				if createdSchedule.Status.Condition == "" {
					return fmt.Errorf("condition not found")
				}

				if createdSchedule.Status.Condition != autoscalingv1.ScheduleAvailable {
					return fmt.Errorf("condition not available: %s", createdSchedule.Status.Condition)
				}

				if existing := metav1.GetControllerOf(&createdSchedule); existing == nil {
					return fmt.Errorf("ownerReference not found")
				}

				return nil
			}, /*timeout*/ time.Second*1 /*pollingInterval*/, time.Millisecond*100).Should(gomega.Succeed())
		})
	})
})
