package calendarboundary

import "time"

// ISOWeek returns the ISO-8601 week-year and week number for a valid Date.
//
// Invalid dates return zero values. ISO weeks start on Monday, and week 1 is
// the week containing January 4.
func ISOWeek(date Date) (year, week int) {
	if !date.Valid() {
		return 0, 0
	}
	weekday := int(date.Weekday())
	isoWeekday := (weekday+6)%7 + 1
	thursdayDays := daysFromCivil(date.Year, int(date.Month), date.Day) + (4 - isoWeekday)
	isoYear, _, _ := civilFromDays(thursdayDays)
	janFourth := daysFromCivil(isoYear, int(time.January), 4)
	janFourthWeekday := (int(weekdayFromDays(janFourth))+6)%7 + 1
	firstMonday := janFourth - (janFourthWeekday - 1)
	weekNumber := (daysFromCivil(date.Year, int(date.Month), date.Day)-firstMonday)/7 + 1
	return isoYear, weekNumber
}
