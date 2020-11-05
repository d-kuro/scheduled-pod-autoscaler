# scheduled-pod-autoscaler

**Work in Progress** GitOps Native Schedule Scaling of Kubernetes Resources.

## Overview

scheduled-pod-autoscaler is made up of two custom resources.

The parent-child relationship can look like this:

```console
$ kubectl tree scheduledpodautoscaler nginx
NAMESPACE  NAME                             READY  REASON  AGE
default    ScheduledPodAutoscaler/nginx     -              6m5s
default    ├─HorizontalPodAutoscaler/nginx  -              6m5s
default    ├─Schedule/test-1                -              6m4s
default    ├─Schedule/test-2                -              6m4s
default    └─Schedule/test-3                -              6m4s
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
NAME    MINPODS   MAXPODS   STATUS      AGE
nginx   3         10        Available   6m52s
```

### Schedule

`Schedule` is a custom resource for defining scheduled scaling.  
You can define multiple children's `Schedule` for the parent `ScheduledPodAutoscaler`.

The `ScheduledPodAutoscaler` controller refers to the `Schedule` and
rewrites `HorizontalPodAutoscaler` created by `ScheduledPodAutoscaler` when it is time for scheduled scaling.
`HorizontalPodAutoscaler` is not managed in Git, so there is no diffs in GitOps.

> 📝 Note: A case of schedule conflicts
>
> In case of a schedule conflict, using the maximum value of min/max replicas.

```console
$ kubectl get schedule -o wide
NAME     REFERENCE   TYPE      STARTTIME          ENDTIME            STARTDAYOFWEEK   ENDDAYOFWEEK   MINPODS   MAXPODS   STATUS      AGE
test-1   nginx       Weekly    20:10              20:15              Saturday         Saturday       1         1         Available   4m49s
test-2   nginx       Daily     20:20              20:25                                              2         2         Available   4m49s
test-3   nginx       OneShot   2020-10-31T20:30   2020-10-31T20:35                                   4         4         Completed   4m49s
```

`Schedule` supports 3 different schedule types.

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

## Spec

### ScheduledPodAutoscaler

| name | type | required | description |
| - | - | - | - |
| `.spec.horizontalPodAutoscalerSpec` | `Object` | required | HorizontalPodAutoscalerSpec is HorizontalPodAutoscaler v2beta2 API spec. ref: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/#horizontalpodautoscaler-v2beta2-autoscaling |

### Schedule

| name | type | required | description |
| - | - | - | - |
| `.spec.scaleTargetRef` | `Object` | required | ScaleTargetRef points to the target resource to scale, and is used to the pods for which metrics should be collected, as well as to actually change the replica count. |
| `.spec.scaleTargetRef.apiVersion` | `string` | optional | API version of the referent. |
| `.spec.scaleTargetRef.kind` | `string` | required | Kind of the referent. |
| `.spec.scaleTargetRef.name` | `string` | required | Name of the referent. |
| `.spec.description` | `string` | optional | Description is schedule description. |
| `.spec.suspend` | `boolean` | optional | Suspend indicates whether to suspend this schedule. |
| `.spec.timeZone` | `string` | optional | TimeZone is the name of the timezone used in the argument of the time.LoadLocation(name string) function. StartTime and EndTime are interpreted as the time in the time zone specified by TimeZone. If not specified, the time will be interpreted as UTC. |
| `.spec.minReplicas` | `integer` | optional | MinReplicas is the lower limit for the number of replicas to which the autoscaler can scale down. It defaults to 1 pod. |
| `.spec.maxReplicas` | `integer` | optional | MaxReplicas is the upper limit for the number of replicas to which the autoscaler can scale up. |
| `.spec.type` | `string` | required | ScheduleType is a type of schedule represented by "Weekly","Daily","OneShot". |
| `.spec.startDayOfWeek` | `string` | optional | StartDayOfWeek is scaling start day of week. Represented by "Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday". |
| `.spec.endDayOfWeek` | `string` | optional | EndDayOfWeek is scaling end day of week. Represented by "Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday". |
| `.spec.startTime` | `string` | required | StartTime is scaling start time. Defined in RFC3339 based format. Different formats are evaluated depending on ScheduleType. e.g. OneShot(yyyy-MM-ddTHH:mm), Weekly(HH:mm), Daily(HH:mm) |
| `.spec.endTime` | `string` | required | EndTime is scaling end time. Defined in RFC3339 based format. Different formats are evaluated depending on ScheduleType. e.g. OneShot(yyyy-MM-ddTHH:mm), Weekly(HH:mm), Daily(HH:mm) |
