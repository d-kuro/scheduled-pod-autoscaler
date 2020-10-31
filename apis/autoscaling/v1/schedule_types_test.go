package v1

import (
	"testing"
	"time"
)

func TestScheduleIsCompleted(t *testing.T) {
	tests := []struct {
		name     string
		spec     ScheduleSpec
		now      time.Time
		expected bool
	}{
		{
			name:     "case[1]",
			spec:     ScheduleSpec{ScheduleType: OneShot, StartTime: "2018-09-01T10:00", EndTime: "2018-09-10T19:00"},
			now:      time.Date(2018, 9, 10, 20, 00, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "case[2]",
			spec:     ScheduleSpec{ScheduleType: OneShot, StartTime: "2018-09-01T10:00", EndTime: "2018-09-10T19:00"},
			now:      time.Date(2018, 9, 10, 18, 59, 0, 0, time.UTC),
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			completed, err := tt.spec.IsCompleted(tt.now)
			if err != nil {
				t.Error(err)

				return
			}

			if completed != tt.expected {
				t.Errorf("%s is not expected condition. actual:%t, expected:%t, endTime: %s",
					tt.now, completed, tt.expected, tt.spec.EndTime)
			}
		})
	}
}
