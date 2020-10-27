package v1

import (
	"testing"
	"time"
)

func TestScheduleSpecContainsDaily(t *testing.T) {
	tests := []struct {
		name     string
		spec     ScheduleSpec
		now      time.Time
		expected bool
	}{
		{
			name:     "case[1]",
			spec:     ScheduleSpec{ScheduleType: TypeDaily, StartTime: "10:00", EndTime: "19:00"},
			now:      time.Date(2018, 9, 2, 11, 00, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "case[2]",
			spec:     ScheduleSpec{ScheduleType: TypeDaily, StartTime: "10:00", EndTime: "19:00"},
			now:      time.Date(2018, 9, 2, 20, 00, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "date changes case[1]",
			spec:     ScheduleSpec{ScheduleType: TypeDaily, StartTime: "23:00", EndTime: "03:00"},
			now:      time.Date(2018, 9, 2, 02, 00, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "date changes case[2]",
			spec:     ScheduleSpec{ScheduleType: TypeDaily, StartTime: "23:00", EndTime: "01:00"},
			now:      time.Date(2018, 9, 2, 02, 00, 0, 0, time.UTC),
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			contains, err := tt.spec.Contains(tt.now)
			if err != nil {
				t.Error(err)

				return
			}

			if contains != tt.expected {
				t.Errorf("%s is not expected condition. actual:%t expected:%t time: %s - %s",
					tt.now, contains, tt.expected,
					tt.spec.StartTime, tt.spec.EndTime)
			}
		})
	}
}

func TestScheduleSpecContainsMonthly(t *testing.T) {
	tests := []struct {
		name     string
		spec     ScheduleSpec
		now      time.Time
		expected bool
	}{
		{
			name:     "case[1]",
			spec:     ScheduleSpec{ScheduleType: TypeMonthly, StartTime: "09-01T10:00", EndTime: "09-10T19:00"},
			now:      time.Date(2018, 9, 10, 11, 00, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "case[2]",
			spec:     ScheduleSpec{ScheduleType: TypeMonthly, StartTime: "09-01T10:00", EndTime: "09-05T19:00"},
			now:      time.Date(2018, 9, 10, 20, 00, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "year changes case[1]",
			spec:     ScheduleSpec{ScheduleType: TypeMonthly, StartTime: "12-01T10:00", EndTime: "01-01T10:00"},
			now:      time.Date(2018, 1, 1, 02, 00, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "year changes case[2]",
			spec:     ScheduleSpec{ScheduleType: TypeMonthly, StartTime: "12-01T10:00", EndTime: "01-01T10:00"},
			now:      time.Date(2018, 1, 2, 02, 00, 0, 0, time.UTC),
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			contains, err := tt.spec.Contains(tt.now)
			if err != nil {
				t.Error(err)

				return
			}

			if contains != tt.expected {
				t.Errorf("%s is not expected condition. actual:%t expected:%t time: %s - %s",
					tt.now, contains, tt.expected,
					tt.spec.StartTime, tt.spec.EndTime)
			}
		})
	}
}

func TestScheduleSpecContainsOneShot(t *testing.T) {
	tests := []struct {
		name     string
		spec     ScheduleSpec
		now      time.Time
		expected bool
	}{
		{
			name:     "case[1]",
			spec:     ScheduleSpec{ScheduleType: TypeOneShot, StartTime: "2018-09-01T10:00", EndTime: "2018-09-10T19:00"},
			now:      time.Date(2018, 9, 10, 11, 00, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "case[2]",
			spec:     ScheduleSpec{ScheduleType: TypeOneShot, StartTime: "2018-09-01T10:00", EndTime: "2018-09-05T19:00"},
			now:      time.Date(2018, 9, 10, 20, 00, 0, 0, time.UTC),
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			contains, err := tt.spec.Contains(tt.now)
			if err != nil {
				t.Error(err)

				return
			}

			if contains != tt.expected {
				t.Errorf("%s is not expected condition. actual:%t expected:%t time: %s - %s",
					tt.now, contains, tt.expected,
					tt.spec.StartTime, tt.spec.EndTime)
			}
		})
	}
}
