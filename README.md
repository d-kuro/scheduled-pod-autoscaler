# scheduled-pod-autoscaler

**Work in Progress** GitOps Native Schedule Scaling of Kubernetes Resources.

## Overview

scheduled-pod-autoscaler is made up of two custom resources.

The parent-child relationship can look like this:

```console
$ kubectl tree scheduledpodautoscaler api
NAMESPACE  NAME                                 READY  REASON  AGE
default    ScheduledPodAutoscaler/api           -              7m20s
default    ├─HorizontalPodAutoscaler/api        -              7m18s
default    ├─Schedule/api-push-notification-01  -              7m20s
default    ├─Schedule/api-push-notification-02  -              7m20s
default    └─Schedule/api-push-notification-03  -              77s
```

### ScheduledPodAutoscaler

ScheduledPodAutoscaler is a custom resource that wraps HorizontalPodAutoscaler.  
The ScheduledPodAutoscaler Controller generates a HorizontalPodAutoscaler from this resource.

The specs of the HorizontalPodAutoscaler defined here will be used when no scheduled scaling is taking place.

```yaml
apiVersion: autoscaling.d-kuro.github.io/v1
kind: ScheduledPodAutoscaler
metadata:
  name: api
spec:
  horizontalPodAutoscalerSpec:
    scaleTargetRef:
      apiVersion: apps/v1
      kind: Deployment
      name: api
    minReplicas: 3
    maxReplicas: 10
    metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
```

```console
$ kubectl get spa # You can use spa as a short name of scheduledpodautoscaler.
NAME   MINPODS   MAXPODS   STATUS      AGE
api    3         10        Available   2m43s
```

### Schedule

Schedule is a custom resource for defining scheduled scaling.  
You can define multiple children's Schedules for the parent ScheduledPodAutoscaler.

Schedule Controller rewrites HorizontalPodAutoscaler created by ScheduledPodAutoscaler when it is time for scheduled scaling
HorizontalPodAutoscaler is not maintained in Git, so there is no diffs in GitOps.

> 📝 Note
>
> If you have a schedule with conflicting scheduled scaling times, the schedule with an earlier start time takes precedence.

```yaml
apiVersion: autoscaling.d-kuro.github.io/v1
kind: Schedule
metadata:
  name: api-push-notification
spec:
  suspend: false
  scaleTargetRef:
    apiVersion: autoscaling.d-kuro.github.io/v1
    kind: ScheduledPodAutoscaler
    name: api
  description: "Scheduled scaling for push notification every day."
  minReplicas: 10
  maxReplicas: 20
  startDayOfWeek: Sunday
  endDayOfWeek: Saturday
  startTime: "18:50"
  endTime: "19:20"
  timeZoneName: Asia/Tokyo
```

```console
$ kubectl get schedule -o wide
NAME                       REFERENCE   STARTTIME   STARTDAYOFWEEK   ENDTIME   ENDDAYOFWEEK   MINPODS   MAXPODS   STATUS        AGE
api-push-notification-01   api         18:50       Sunday           19:20     Saturday       10        20        Available     6m8s
api-push-notification-02   api         21:40       Tuesday          22:00     Tuesday        5         15        Progressing   6m8s
api-push-notification-03   api         11:50       Monday           13:00     Monday         5         15        Suspend       5s
```

## Install

> TBD
