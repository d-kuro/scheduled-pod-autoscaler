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

`ScheduledPodAutoscaler` is a custom resource that wraps `HorizontalPodAutoscaler`.  
The `ScheduledPodAutoscaler` Controller generates a `HorizontalPodAutoscaler` from this resource.

The specs of the `HorizontalPodAutoscaler` defined here will be used when no scheduled scaling is taking place.

for example:

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

`Schedule` is a custom resource for defining scheduled scaling.  
You can define multiple children's `Schedule` for the parent `ScheduledPodAutoscaler`.

The `ScheduledPodAutoscaler` controller refers to the `Schedule` and
rewrites `HorizontalPodAutoscaler` created by `ScheduledPodAutoscaler` when it is time for scheduled scaling.
`HorizontalPodAutoscaler` is not managed in Git, so there is no diffs in GitOps.

`Schedule` supports 4 different schedule types.

#### type: Monthly

Write the time in the format of `MM-ddTHH:mm`.

```yaml
apiVersion: autoscaling.d-kuro.github.io/v1
kind: Schedule
metadata:
  name: api-push-notification
spec:
  scaleTargetRef:
    apiVersion: autoscaling.d-kuro.github.io/v1
    kind: ScheduledPodAutoscaler
    name: api
  minReplicas: 10
  maxReplicas: 20
  type: Monthly
  startTime: "09-01T09:50"
  endTime: "09-10T19:00"
  timeZone: Asia/Tokyo
```

#### type: Weekly

Write the time in the format of `HH:mm` and specify the day of the week.

```yaml
apiVersion: autoscaling.d-kuro.github.io/v1
kind: Schedule
metadata:
  name: api-push-notification
spec:
  scaleTargetRef:
    apiVersion: autoscaling.d-kuro.github.io/v1
    kind: ScheduledPodAutoscaler
    name: api
  minReplicas: 10
  maxReplicas: 20
  type: Weekly
  startDayOfWeek: Monday
  startTime: "11:50"
  endDayOfWeek: Wednesday
  endTime: "13:00"
  timeZone: Asia/Tokyo
```

#### type: Daily

Write the time in the format of `HH:mm`.

```yaml
apiVersion: autoscaling.d-kuro.github.io/v1
kind: Schedule
metadata:
  name: api-push-notification
spec:
  scaleTargetRef:
    apiVersion: autoscaling.d-kuro.github.io/v1
    kind: ScheduledPodAutoscaler
    name: api
  minReplicas: 10
  maxReplicas: 20
  type: Daily
  startTime: "11:50"
  endTime: "13:00"
  timeZone: Asia/Tokyo
```

#### type: OneShot

Write the time in the format of `yyyy-MM-ddTHH:mm`.

```yaml
apiVersion: autoscaling.d-kuro.github.io/v1
kind: Schedule
metadata:
  name: api-push-notification
spec:
  scaleTargetRef:
    apiVersion: autoscaling.d-kuro.github.io/v1
    kind: ScheduledPodAutoscaler
    name: api
  minReplicas: 10
  maxReplicas: 20
  type: OneShot
  startTime: "2020-09-01T10:00"
  endTime: "2020-09-10T19:00"
  timeZone: Asia/Tokyo
```

## Install

> TBD
