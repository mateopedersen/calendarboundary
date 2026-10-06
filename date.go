package calendarboundary

import (
	"errors"
	"fmt"
	"time"
)

const (
	minYear = 1
	maxYear = 9999
)

var (
	// ErrInvalidDate indicates that a civil date is outside the supported range.
	ErrInvalidDate = errors.New("calendarboundary: invalid civil date")
	// ErrDateOutOfRange indicates that date arithmetic would leave years 1–9999.
	ErrDateOutOfRange = errors.New("calendarboundary: date is outside years 1–9999")
)

// Date is a Gregorian civil date without a time or time zone.
//
// A Date is valid when Year is 1–9999, Month is January–December, and Day
// exists in that month. Use NewDate to validate a value at an API boundary.
type Date struct {
	Year  int
	Month time.Month
	Day   int
}

// NewDate returns a validated Gregorian civil date.
func NewDate(year int, month time.Month, day int) (Date, error) {
	d := Date{Year: year, Month: month, Day: day}
	if !d.Valid() {
		return Date{}, fmt.Errorf("%w: %04d-%02d-%02d", ErrInvalidDate, year, month, day)
	}
	return d, nil
}

// Valid reports whether d is a supported civil date.
func (d Date) Valid() bool {
	return d.Year >= minYear && d.Year <= maxYear &&
		d.Month >= time.January && d.Month <= time.December &&
		d.Day >= 1 && d.Day <= DaysInMonth(d.Year, d.Month)
}

// String formats d as an ISO-style YYYY-MM-DD date.
func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}

// Weekday returns the weekday of d, or time.Sunday for an invalid Date.
func (d Date) Weekday() time.Weekday {
	if !d.Valid() {
		return time.Sunday
	}
	return weekdayFromDays(daysFromCivil(d.Year, int(d.Month), d.Day))
}

// ISOWeek returns the ISO-8601 week-year and week number for d.
//
// The week-year may differ from d.Year near New Year's Day. An invalid Date
// returns zero values.
func (d Date) ISOWeek() (year, week int) {
	if !d.Valid() {
		return 0, 0
	}
	return ISOWeek(d)
}

// AddDays returns d shifted by offset Gregorian days.
func (d Date) AddDays(offset int) (Date, error) {
	if !d.Valid() {
		return Date{}, ErrInvalidDate
	}
	// No valid result can be reached using a larger offset from this range.
	const maxSupportedDays = 366*10000 + 7
	if offset > maxSupportedDays || offset < -maxSupportedDays {
		return Date{}, ErrDateOutOfRange
	}
	year, month, day := civilFromDays(daysFromCivil(d.Year, int(d.Month), d.Day) + offset)
	result := Date{Year: year, Month: time.Month(month), Day: day}
	if !result.Valid() {
		return Date{}, ErrDateOutOfRange
	}
	return result, nil
}

// IsLeapYear reports whether year follows the Gregorian leap-year rule.
func IsLeapYear(year int) bool {
	return year%400 == 0 || (year%4 == 0 && year%100 != 0)
}

// DaysInMonth returns a month's Gregorian length, or zero for an invalid year
// or month. Supported civil years are 1 through 9999.
func DaysInMonth(year int, month time.Month) int {
	if year < minYear || year > maxYear || month < time.January || month > time.December {
		return 0
	}
	switch month {
	case time.April, time.June, time.September, time.November:
		return 30
	case time.February:
		if IsLeapYear(year) {
			return 29
		}
		return 28
	default:
		return 31
	}
}

func weekdayFromDays(days int) time.Weekday {
	// daysFromCivil is relative to an arbitrary origin; 1970-01-01 is day 0,
	// a Thursday, so the Sunday-zero index is (days+4) mod 7.
	weekday := (days + 4) % 7
	if weekday < 0 {
		weekday += 7
	}
	return time.Weekday(weekday)
}

// daysFromCivil and civilFromDays are inverse proleptic Gregorian transforms.
// The day origin is the Unix epoch; both use floor division for negative years.
func daysFromCivil(year, month, day int) int {
	y := year
	if month <= 2 {
		y--
	}
	era := floorDiv(y, 400)
	yoe := y - era*400
	mp := month
	if month > 2 {
		mp -= 3
	} else {
		mp += 9
	}
	doy := (153*mp+2)/5 + day - 1
	doe := yoe*365 + yoe/4 - yoe/100 + doy
	return era*146097 + doe - 719468
}

func civilFromDays(days int) (year, month, day int) {
	z := days + 719468
	era := floorDiv(z, 146097)
	doe := z - era*146097
	yoe := (doe - doe/1460 + doe/36524 - doe/146096) / 365
	y := yoe + era*400
	doy := doe - (365*yoe + yoe/4 - yoe/100)
	mp := (5*doy + 2) / 153
	d := doy - (153*mp+2)/5 + 1
	m := mp + 3
	if mp >= 10 {
		m = mp - 9
	}
	if m <= 2 {
		y++
	}
	return y, m, d
}

func floorDiv(a, b int) int {
	q := a / b
	r := a % b
	if r != 0 && (r < 0) != (b < 0) {
		q--
	}
	return q
}
