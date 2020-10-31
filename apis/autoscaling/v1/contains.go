package v1

import (
	"fmt"
	"time"

	// embed tzdata.
	_ "time/tzdata"
)

var weekdays = map[string]time.Weekday{
	"Sunday":    time.Weekday(0),
	"Monday":    time.Weekday(1),
	"Tuesday":   time.Weekday(2),
	"Wednesday": time.Weekday(3),
	"Thursday":  time.Weekday(4),
	"Friday":    time.Weekday(5),
	"Saturday":  time.Weekday(6),
}

func (s *ScheduleSpec) Contains(now time.Time) (bool, error) {
	location, err := time.LoadLocation(s.TimeZone)
	if err != nil {
		return false, fmt.Errorf("failed to load location %s: %w", s.TimeZone, err)
	}

	now = now.In(location)

	switch s.ScheduleType {
	case Daily:
		return s.containsDaily(now, location)
	case Weekly:
		return s.containsWeekly(now, location)
	case OneShot:
		return s.containsOneShot(now, location)
	default:
		return false, fmt.Errorf("unsupported schedule types: %s", s.ScheduleType)
	}
}

func (s *ScheduleSpec) containsDaily(now time.Time, location *time.Location) (bool, error) {
	startTime, endTime, err := s.normalizeTime(now, location)
	if err != nil {
		return false, err
	}

	// true if now is [startTime, endTime)
	return (now.Equal(startTime) || now.After(startTime)) && now.Before(endTime), nil
}

func (s *ScheduleSpec) containsWeekly(now time.Time, location *time.Location) (bool, error) {
	startTime, endTime, err := s.normalizeTime(now, location)
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

func (s *ScheduleSpec) containsOneShot(now time.Time, location *time.Location) (bool, error) {
	startTime, err := time.ParseInLocation("2006-01-02T15:04", s.StartTime, location)
	if err != nil {
		return false, fmt.Errorf("startTime cannot be parsed: %w", err)
	}

	endTime, err := time.ParseInLocation("2006-01-02T15:04", s.EndTime, location)
	if err != nil {
		return false, fmt.Errorf("endTime cannot be parsed: %w", err)
	}

	// true if now is [startTime, endTime)
	return (now.Equal(startTime) || now.After(startTime)) && now.Before(endTime), nil
}

func (s *ScheduleSpec) normalizeDateTime(now time.Time, location *time.Location) (normalizedStartTime time.Time, normalizedEndTime time.Time, err error) {
	startTime, err := time.ParseInLocation("02T15:04", s.StartTime, location)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("startTime cannot be parsed: %w", err)
	}

	endTime, err := time.ParseInLocation("02T15:04", s.EndTime, location)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("endTime cannot be parsed: %w", err)
	}

	normalizedStartTime = time.Date(
		now.Year(), now.Month(), startTime.Day(),
		startTime.Hour(), startTime.Minute(), 0, 0, location)

	normalizedEndTime = time.Date(
		now.Year(), now.Month(), endTime.Day(),
		endTime.Hour(), endTime.Minute(), 0, 0, location)

	if normalizedEndTime.Before(normalizedStartTime) {
		if now.Before(normalizedStartTime) {
			normalizedStartTime = normalizedStartTime.AddDate(-1, 0, 0)
		} else {
			normalizedEndTime = normalizedEndTime.AddDate(1, 0, 0)
		}
	}

	return normalizedStartTime, normalizedEndTime, nil
}

func (s *ScheduleSpec) normalizeTime(now time.Time, location *time.Location) (normalizedStartTime time.Time, normalizedEndTime time.Time, err error) {
	startTime, err := time.ParseInLocation("15:04", s.StartTime, location)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("startTime cannot be parsed: %w", err)
	}

	endTime, err := time.ParseInLocation("15:04", s.EndTime, location)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("endTime cannot be parsed: %w", err)
	}

	normalizedStartTime = time.Date(
		now.Year(), now.Month(), now.Day(),
		startTime.Hour(), startTime.Minute(), 0, 0, location)

	normalizedEndTime = time.Date(
		now.Year(), now.Month(), now.Day(),
		endTime.Hour(), endTime.Minute(), 0, 0, location)

	if normalizedEndTime.Before(normalizedStartTime) {
		if now.Before(normalizedStartTime) {
			normalizedStartTime = normalizedStartTime.AddDate(0, 0, -1)
		} else {
			normalizedEndTime = normalizedEndTime.AddDate(0, 0, 1)
		}
	}

	return normalizedStartTime, normalizedEndTime, nil
}

func (s *ScheduleSpec) normalizeWeekday(startTime time.Time) (
	weekdayToday time.Weekday, startWeekDay time.Weekday, endWeekDay time.Weekday, err error) {
	startWeekDay, found := weekdays[s.StartDayOfWeek]
	if !found {
		return 0, 0, 0, fmt.Errorf("start-day-of-week %s is not found", s.StartDayOfWeek)
	}

	endWeekDay, found = weekdays[s.EndDayOfWeek]
	if !found {
		return 0, 0, 0, fmt.Errorf("end-day-of-week %s is invalid", s.EndDayOfWeek)
	}

	weekdayToday = startTime.Weekday()
	if startWeekDay > endWeekDay {
		// normalize weekday
		endWeekDay = 7 - startWeekDay + endWeekDay
		weekdayToday = (7 + weekdayToday - startWeekDay) % 7
		startWeekDay = 0
	}

	return weekdayToday, startWeekDay, endWeekDay, nil
}
