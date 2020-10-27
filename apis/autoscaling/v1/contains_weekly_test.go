package v1

import (
	"testing"
	"time"
)

func TestScheduleSpecContainsWeekly(t *testing.T) {
	tests := []struct {
		name     string
		spec     ScheduleSpec
		now      time.Time
		expected bool
	}{
		// 1 day case
		{
			name:     "1day case[1]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Monday", EndDayOfWeek: "Monday", StartTime: "10:00", EndTime: "19:00"},
			now:      time.Date(2018, 9, 2, 11, 00, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "1day case[2]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Monday", EndDayOfWeek: "Monday", StartTime: "10:00", EndTime: "19:00"},
			now:      time.Date(2018, 9, 3, 9, 00, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "1day case[3]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Monday", EndDayOfWeek: "Monday", StartTime: "10:00", EndTime: "19:00"},
			now:      time.Date(2018, 9, 3, 10, 00, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "1day case[4]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Monday", EndDayOfWeek: "Monday", StartTime: "10:00", EndTime: "19:00"},
			now:      time.Date(2018, 9, 3, 11, 00, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "1day case[5]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Monday", EndDayOfWeek: "Monday", StartTime: "10:00", EndTime: "19:00"},
			now:      time.Date(2018, 9, 3, 18, 00, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "1day case[6]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Monday", EndDayOfWeek: "Monday", StartTime: "10:00", EndTime: "19:00"},
			now:      time.Date(2018, 9, 3, 19, 00, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "1day case[7]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Monday", EndDayOfWeek: "Monday", StartTime: "10:00", EndTime: "19:00"},
			now:      time.Date(2018, 9, 3, 20, 00, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "1day case[8]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Monday", EndDayOfWeek: "Monday", StartTime: "10:00", EndTime: "19:00"},
			now:      time.Date(2018, 9, 4, 11, 00, 0, 0, time.UTC),
			expected: false,
		},

		// 1 day case with minutes
		{
			name:     "1day case with minutes[1]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Monday", EndDayOfWeek: "Monday", StartTime: "10:00", EndTime: "10:30"},
			now:      time.Date(2018, 9, 3, 9, 50, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "1day case with minutes[2]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Monday", EndDayOfWeek: "Monday", StartTime: "10:00", EndTime: "10:30"},
			now:      time.Date(2018, 9, 3, 10, 20, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "1day case with minutes[3]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Monday", EndDayOfWeek: "Monday", StartTime: "10:00", EndTime: "10:30"},
			now:      time.Date(2018, 9, 3, 9, 50, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "1day case with minutes[4]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Monday", EndDayOfWeek: "Monday", StartTime: "10:00", EndTime: "10:30"},
			now:      time.Date(2018, 9, 3, 10, 40, 0, 0, time.UTC),
			expected: false,
		},

		// EndTime is over 24 o'clock.
		{
			name:     "EndTime is over 24o'clock[1]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Monday", EndDayOfWeek: "Monday", StartTime: "15:00", EndTime: "01:00"},
			now:      time.Date(2018, 9, 3, 0, 00, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "EndTime is over 24o'clock[2]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Monday", EndDayOfWeek: "Monday", StartTime: "15:00", EndTime: "01:00"},
			now:      time.Date(2018, 9, 3, 14, 00, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "EndTime is over 24o'clock[3]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Monday", EndDayOfWeek: "Monday", StartTime: "15:00", EndTime: "01:00"},
			now:      time.Date(2018, 9, 3, 15, 00, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "EndTime is over 24o'clock[4]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Monday", EndDayOfWeek: "Monday", StartTime: "15:00", EndTime: "01:00"},
			now:      time.Date(2018, 9, 3, 19, 00, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "EndTime is over 24o'clock[5]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Monday", EndDayOfWeek: "Monday", StartTime: "15:00", EndTime: "01:00"},
			now:      time.Date(2018, 9, 4, 0, 00, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "EndTime is over 24o'clock[6]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Monday", EndDayOfWeek: "Monday", StartTime: "15:00", EndTime: "01:00"},
			now:      time.Date(2018, 9, 4, 1, 00, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "EndTime is over 24o'clock[7]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Monday", EndDayOfWeek: "Monday", StartTime: "15:00", EndTime: "01:00"},
			now:      time.Date(2018, 9, 4, 2, 00, 0, 0, time.UTC),
			expected: false,
		},

		// EndTime is over 24 o'clock with minutes case.
		{
			name:     "EndTime is over 24o'clock with minutes case[1]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Monday", EndDayOfWeek: "Monday", StartTime: "15:00", EndTime: "01:30"},
			now:      time.Date(2018, 9, 4, 1, 20, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "EndTime is over 24o'clock with minutes case[2]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Monday", EndDayOfWeek: "Monday", StartTime: "15:00", EndTime: "01:30"},
			now:      time.Date(2018, 9, 4, 1, 30, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "EndTime is over 24o'clock with minutes case[3]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Monday", EndDayOfWeek: "Monday", StartTime: "15:00", EndTime: "01:30"},
			now:      time.Date(2018, 9, 5, 1, 20, 0, 0, time.UTC),
			expected: false,
		},

		// EndTime is over 24 o'clock, 2 days.
		{
			name:     "EndTime is over 24o'clock 2days[1]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Monday", EndDayOfWeek: "Tuesday", StartTime: "15:00", EndTime: "01:30"},
			now:      time.Date(2018, 9, 3, 16, 00, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "EndTime is over 24o'clock 2days[2]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Monday", EndDayOfWeek: "Tuesday", StartTime: "15:00", EndTime: "01:30"},
			now:      time.Date(2018, 9, 4, 16, 00, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "EndTime is over 24o'clock 2days[3]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Monday", EndDayOfWeek: "Tuesday", StartTime: "15:00", EndTime: "01:30"},
			now:      time.Date(2018, 9, 5, 1, 00, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "EndTime is over 24o'clock 2days[4]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Monday", EndDayOfWeek: "Tuesday", StartTime: "15:00", EndTime: "01:30"},
			now:      time.Date(2018, 9, 5, 2, 00, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "EndTime is over 24o'clock 2days[5]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Monday", EndDayOfWeek: "Tuesday", StartTime: "15:00", EndTime: "01:30"},
			now:      time.Date(2018, 9, 6, 16, 00, 0, 0, time.UTC),
			expected: false,
		},

		// the middle day in 3 days
		{
			name:     "the middle day in 3days[1]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Monday", EndDayOfWeek: "Wednesday", StartTime: "15:00", EndTime: "01:30"},
			now:      time.Date(2018, 9, 4, 14, 00, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "the middle day in 3days[2]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Monday", EndDayOfWeek: "Wednesday", StartTime: "15:00", EndTime: "01:30"},
			now:      time.Date(2018, 9, 4, 15, 00, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "the middle day in 3days[3]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Monday", EndDayOfWeek: "Wednesday", StartTime: "15:00", EndTime: "01:30"},
			now:      time.Date(2018, 9, 4, 16, 00, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "the middle day in 3days[4]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Monday", EndDayOfWeek: "Wednesday", StartTime: "15:00", EndTime: "01:30"},
			now:      time.Date(2018, 9, 4, 23, 00, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "the middle day in 3days[5]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Monday", EndDayOfWeek: "Wednesday", StartTime: "15:00", EndTime: "01:30"},
			now:      time.Date(2018, 9, 5, 0, 00, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "the middle day in 3days[6]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Monday", EndDayOfWeek: "Wednesday", StartTime: "15:00", EndTime: "01:30"},
			now:      time.Date(2018, 9, 5, 1, 00, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "the middle day in 3days[7]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Monday", EndDayOfWeek: "Wednesday", StartTime: "15:00", EndTime: "01:30"},
			now:      time.Date(2018, 9, 5, 2, 00, 0, 0, time.UTC),
			expected: false,
		},

		// StartDayOfWeek is bigger than EndDayOfWeek
		{
			name:     "StartDayOfWeek is bigger than EndDayOfWeek[1]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Tuesday", EndDayOfWeek: "Sunday", StartTime: "15:00", EndTime: "01:30"},
			now:      time.Date(2018, 9, 3, 1, 20, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "StartDayOfWeek is bigger than EndDayOfWeek[2]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Tuesday", EndDayOfWeek: "Sunday", StartTime: "15:00", EndTime: "01:30"},
			now:      time.Date(2018, 9, 3, 16, 20, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "StartDayOfWeek is bigger than EndDayOfWeek[3]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Tuesday", EndDayOfWeek: "Sunday", StartTime: "15:00", EndTime: "01:30"},
			now:      time.Date(2018, 9, 4, 1, 20, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "StartDayOfWeek is bigger than EndDayOfWeek[4]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Tuesday", EndDayOfWeek: "Sunday", StartTime: "15:00", EndTime: "01:30"},
			now:      time.Date(2018, 9, 5, 1, 20, 0, 0, time.UTC),
			expected: true,
		},

		// everyday
		{
			name:     "everyday[1]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Sunday", EndDayOfWeek: "Saturday", StartTime: "15:00", EndTime: "01:30"},
			now:      time.Date(2018, 9, 3, 1, 20, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "everyday[2]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Sunday", EndDayOfWeek: "Saturday", StartTime: "15:00", EndTime: "01:30"},
			now:      time.Date(2018, 9, 3, 16, 00, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "everyday[3]",
			spec:     ScheduleSpec{ScheduleType: TypeWeekly, StartDayOfWeek: "Sunday", EndDayOfWeek: "Saturday", StartTime: "15:00", EndTime: "01:30"},
			now:      time.Date(2018, 9, 5, 1, 20, 0, 0, time.UTC),
			expected: true,
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
				t.Errorf("%s is not expected condition. actual:%t expected:%t time: %s - %s day: %s - %s",
					tt.now, contains, tt.expected,
					tt.spec.StartTime, tt.spec.EndTime,
					tt.spec.StartDayOfWeek, tt.spec.EndDayOfWeek)
				t.Errorf("%s is %s", tt.now, tt.now.Weekday())
			}
		})
	}
}
