package calendarboundary

import (
	"errors"
	"testing"
	"time"
)

func TestGregorianLeapYearsAndMonthLengths(t *testing.T) {
	tests := []struct {
		year int
		leap bool
	}{
		{year: 2024, leap: true},
		{year: 2028, leap: true},
		{year: 1900, leap: false},
		{year: 2000, leap: true},
		{year: 2100, leap: false},
	}
	for _, test := range tests {
		if got := IsLeapYear(test.year); got != test.leap {
			t.Errorf("IsLeapYear(%d) = %v, want %v", test.year, got, test.leap)
		}
	}
	for _, test := range []struct {
		year  int
		month time.Month
		want  int
	}{
		{2027, time.February, 28},
		{2028, time.February, 29},
		{2027, time.April, 30},
		{2027, time.January, 31},
	} {
		if got := DaysInMonth(test.year, test.month); got != test.want {
			t.Errorf("DaysInMonth(%d, %s) = %d, want %d", test.year, test.month, got, test.want)
		}
	}
}

func TestDateValidationAndYearBoundary(t *testing.T) {
	if _, err := NewDate(1900, time.February, 29); !errors.Is(err, ErrInvalidDate) {
		t.Fatalf("NewDate for non-leap February 29 error = %v, want ErrInvalidDate", err)
	}
	leapDay, err := NewDate(2000, time.February, 29)
	if err != nil {
		t.Fatalf("NewDate leap day: %v", err)
	}
	if leapDay.String() != "2000-02-29" {
		t.Errorf("String() = %q, want 2000-02-29", leapDay.String())
	}
	end, _ := NewDate(2026, time.December, 31)
	got, err := end.AddDays(1)
	if err != nil {
		t.Fatalf("AddDays year boundary: %v", err)
	}
	if got.String() != "2027-01-01" {
		t.Errorf("AddDays(1) = %s, want 2027-01-01", got)
	}
	start, _ := NewDate(2027, time.January, 1)
	got, err = start.AddDays(-1)
	if err != nil || got.String() != "2026-12-31" {
		t.Errorf("AddDays(-1) = %s, %v; want 2026-12-31, nil", got, err)
	}
	last, _ := NewDate(9999, time.December, 31)
	if _, err := last.AddDays(1); !errors.Is(err, ErrDateOutOfRange) {
		t.Errorf("AddDays beyond supported years error = %v, want ErrDateOutOfRange", err)
	}
}

func TestWeekdayAndISOWeekEdges(t *testing.T) {
	weekday, _ := NewDate(2027, time.January, 1)
	if got := weekday.Weekday(); got != time.Friday {
		t.Errorf("2027-01-01 weekday = %s, want Friday", got)
	}
	for _, test := range []struct {
		date Date
		year int
		week int
	}{
		{Date{2021, time.January, 1}, 2020, 53},
		{Date{2021, time.January, 4}, 2021, 1},
		{Date{2020, time.December, 31}, 2020, 53},
		{Date{2016, time.January, 1}, 2015, 53},
		{Date{2018, time.December, 31}, 2019, 1},
	} {
		year, week := ISOWeek(test.date)
		if year != test.year || week != test.week {
			t.Errorf("ISOWeek(%s) = %d-W%02d, want %d-W%02d", test.date, year, week, test.year, test.week)
		}
	}
}

func TestMonthGridShapeWeekStartsAndEvents(t *testing.T) {
	eventDate := Date{2028, time.February, 29}
	events := map[Date][]string{eventDate: {"leap-day"}}
	grid, err := NewMonthGrid(2028, time.February, Options{
		WeekStart:       time.Monday,
		WeekNumbers:     true,
		ShowOutsideDays: true,
		WeekdayLabels:   []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"},
		Events:          events,
	})
	if err != nil {
		t.Fatalf("NewMonthGrid: %v", err)
	}
	if len(grid.Weeks) != 5 || !grid.ShowWeekNumbers {
		t.Fatalf("grid shape = %d rows, week numbers %v; want 5 and true", len(grid.Weeks), grid.ShowWeekNumbers)
	}
	if grid.WeekdayLabels[0] != "Mon" || grid.WeekdayLabels[6] != "Sun" {
		t.Errorf("Monday-first labels = %v", grid.WeekdayLabels)
	}
	last := grid.Weeks[4].Cells[1]
	if !last.InMonth || last.Date != eventDate || !last.ShowDayNumber || len(last.Events) != 1 || last.Events[0] != "leap-day" {
		t.Errorf("February 29 cell = %+v", last)
	}
	if last.ISOYear != 2028 {
		t.Errorf("ISO year for February 29 = %d, want 2028", last.ISOYear)
	}

	sixRows, err := NewMonthGrid(2021, time.May, Options{WeekStart: time.Sunday, ShowOutsideDays: true})
	if err != nil {
		t.Fatalf("NewMonthGrid six rows: %v", err)
	}
	if len(sixRows.Weeks) != 6 {
		t.Errorf("May 2021 row count = %d, want 6", len(sixRows.Weeks))
	}
	first := sixRows.Weeks[0].Cells[0]
	if !first.HasDate || first.InMonth || !first.ShowDayNumber || first.Date.String() != "2021-04-25" {
		t.Errorf("outside leading cell = %+v", first)
	}

	fourWeeksMonth, err := NewMonthGrid(2021, time.February, Options{WeekStart: time.Monday})
	if err != nil {
		t.Fatalf("NewMonthGrid short month: %v", err)
	}
	if len(fourWeeksMonth.Weeks) != 5 {
		t.Errorf("adaptive minimum row count = %d, want 5", len(fourWeeksMonth.Weeks))
	}
}

func TestFixedRowsInvalidAndBlankCells(t *testing.T) {
	if _, err := NewMonthGrid(2021, time.May, Options{FixedRows: 5}); err == nil {
		t.Fatal("NewMonthGrid accepted a five-row grid for a six-row month")
	}
	grid, err := NewMonthGrid(2028, time.February, Options{WeekStart: time.Monday, Blank: true})
	if err != nil {
		t.Fatalf("NewMonthGrid blank: %v", err)
	}
	cell := grid.Weeks[0].Cells[0]
	if !cell.HasDate || cell.InMonth || cell.ShowDayNumber || cell.Day == 0 {
		t.Errorf("blank outside cell = %+v; want actual date without printed number", cell)
	}
	if _, err := NewMonthGrid(2028, time.February, Options{FixedRows: 4}); err == nil {
		t.Fatal("NewMonthGrid accepted unsupported fixed row count")
	}
}

func TestYearGridIncludesAllMonths(t *testing.T) {
	grids, err := NewYearGrid(2028, Options{WeekStart: time.Monday})
	if err != nil {
		t.Fatalf("NewYearGrid: %v", err)
	}
	if grids[0].Month != time.January || grids[11].Month != time.December {
		t.Errorf("year endpoints = %s, %s", grids[0].Month, grids[11].Month)
	}
	if DaysInMonth(grids[1].Year, grids[1].Month) != 29 {
		t.Errorf("year grid February days = %d, want 29", DaysInMonth(grids[1].Year, grids[1].Month))
	}
}
