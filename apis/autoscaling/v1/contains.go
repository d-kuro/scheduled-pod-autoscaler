package v1

import (
	"fmt"
	"time"
)

var weekdays = map[string]time.Weekday{
	"Monday":    time.Weekday(1),
	"Tuesday":   time.Weekday(2),
	"Wednesday": time.Weekday(3),
	"Thursday":  time.Weekday(4),
	"Friday":    time.Weekday(5),
	"Saturday":  time.Weekday(6),
	"Sunday":    time.Weekday(0),
}

func (s *ScheduleSpec) Contains(now time.Time) (bool, error) {
	startTime, endTime, err := s.normalizeTime(now)
	if err != nil {
		return false, err
	}
	weekdayToday, startWeekDay, endWeekDay, err := s.normalizeWeekday(startTime)
	if err != nil {
		return false, err
	}
	if startWeekDay <= weekdayToday && weekdayToday <= endWeekDay {
		// true if now is [startTime, endTime)
		return (now.Equal(startTime) || now.After(startTime)) && now.Before(endTime), nil
	}
	return false, nil
}

func (s *ScheduleSpec) normalizeTime(now time.Time) (time.Time, time.Time, error) {
	startTime, err := time.Parse("15:04", s.StartTime)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("startTime cannot be parsed: %w", err)
	}
	endTime, err := time.Parse("15:04", s.EndTime)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("endTime cannot be parsed: %w", err)
	}
	normalizedStartTime := time.Date(
		now.Year(), now.Month(), now.Day(),
		startTime.Hour(), startTime.Minute(), 0, 0, time.UTC)
	normalizedEndTime := time.Date(
		now.Year(), now.Month(), now.Day(),
		endTime.Hour(), endTime.Minute(), 0, 0, time.UTC)
	if normalizedEndTime.Before(normalizedStartTime) {
		if now.Hour() <= endTime.Hour() && now.Minute() <= endTime.Minute() {
			normalizedStartTime = normalizedStartTime.AddDate(0, 0, -1)
		} else {
			normalizedEndTime = normalizedEndTime.AddDate(0, 0, 1)
		}
	}
	return normalizedStartTime, normalizedEndTime, nil
}

func (s *ScheduleSpec) normalizeWeekday(startTime time.Time) (
	time.Weekday, time.Weekday, time.Weekday, error) {
	startWeekDay, found := weekdays[s.StartDayOfWeek]
	if !found {
		return 0, 0, 0, fmt.Errorf("start-day-of-week %s is not found", s.StartDayOfWeek)
	}

	endWeekDay, found := weekdays[s.EndDayOfWeek]
	if !found {
		return 0, 0, 0, fmt.Errorf("end-day-of-week %s is invalid", s.EndDayOfWeek)
	}

	weekdayToday := startTime.Weekday()
	if startWeekDay > endWeekDay {
		// normalize weekday
		endWeekDay = 7 - startWeekDay + endWeekDay
		weekdayToday = (7 + weekdayToday - startWeekDay) % 7
		startWeekDay = 0
	}
	return weekdayToday, startWeekDay, endWeekDay, nil
}
