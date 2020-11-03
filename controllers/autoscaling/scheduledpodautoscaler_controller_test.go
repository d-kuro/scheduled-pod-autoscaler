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
				name      = "scheduled-pod-autoscaler-test"
				namespace = "default"
			)

			ctx := context.Background()

			spa := &autoscalingv1.ScheduledPodAutoscaler{
				TypeMeta: metav1.TypeMeta{
					APIVersion: autoscalingv1.GroupVersion.String(),
					Kind:       "ScheduledPodAutoscaler",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: autoscalingv1.ScheduledPodAutoscalerSpec{
					HorizontalPodAutoscalerSpec: hpav2beta2.HorizontalPodAutoscalerSpec{
						ScaleTargetRef: hpav2beta2.CrossVersionObjectReference{
							APIVersion: "apps/v1",
							Kind:       "Deployment",
							Name:       name,
						},
						MinReplicas: testutil.ToPointerInt32(1),
						MaxReplicas: 3,
						Metrics: []hpav2beta2.MetricSpec{
							{
								Type: "Resource",
								Resource: &hpav2beta2.ResourceMetricSource{
									Name: "cpu",
									Target: hpav2beta2.MetricTarget{
										Type:               hpav2beta2.MetricTargetType("Utilization"),
										AverageUtilization: testutil.ToPointerInt32(50),
									},
								},
							},
						},
					},
				},
			}

			err := k8sClient.Create(ctx, spa)
			gomega.Expect(err).Should(gomega.Succeed())

			var createdSPA autoscalingv1.ScheduledPodAutoscaler
			var createdHPA hpav2beta2.HorizontalPodAutoscaler

			gomega.Eventually(func() error {
				if err := k8sClient.Get(ctx, client.ObjectKey{Name: name, Namespace: namespace}, &createdHPA); err != nil {
					return err
				}

				if diff := cmp.Diff(spa.Spec.HorizontalPodAutoscalerSpec, createdHPA.Spec); diff != "" {
					return fmt.Errorf("created HPA mismatch (-want +got):\\n%s", diff)
				}

				if err := k8sClient.Get(ctx, client.ObjectKey{Name: name, Namespace: namespace}, &createdSPA); err != nil {
					return err
				}

				if createdSPA.Status.Condition == "" {
					return fmt.Errorf("condition not found")
				}

				if createdSPA.Status.Condition != autoscalingv1.ScheduledPodAutoscalerAvailable {
					return fmt.Errorf("condition not available: %s", createdSPA.Status.Condition)
				}

				return nil
			}, /*timeout*/ time.Second*1 /*pollingInterval*/, time.Millisecond*100).Should(gomega.Succeed())
		})
	})
})
