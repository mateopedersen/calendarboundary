package calendarboundary

import (
	"fmt"
	"time"
)

// Options configures the shape and annotations of a month grid.
//
// The zero value starts weeks on Sunday, generates the minimum five or six
// complete rows, and hides adjacent-month day numbers. WeekdayLabels, when
// non-empty, must contain seven values in Sunday-to-Saturday order. Events are
// keyed by civil date; their labels are copied into matching cells.
type Options struct {
	// WeekStart is the first weekday in a row; its zero value is Sunday.
	WeekStart time.Weekday
	// FixedRows requests five or six rows; zero chooses the adaptive height.
	FixedRows int
	// ShowOutsideDays controls whether adjacent-month day numbers are displayed.
	ShowOutsideDays bool
	// WeekNumbers indicates that consumers should display each row's ISO label.
	WeekNumbers bool
	// Blank suppresses day numbers while preserving dates, events, and grid shape.
	Blank bool
	// WeekdayLabels contains seven Sunday-to-Saturday labels, or is empty for English abbreviations.
	WeekdayLabels []string
	// Events maps civil dates to event labels shown in matching visible cells.
	Events map[Date][]string
}

// Cell is one day position in a month grid.
//
// HasDate is false only for grid positions that fall outside the supported
// range (years 1–9999). ShowDayNumber distinguishes hidden outside dates and
// blank-planning cells from ordinary numbered dates.
type Cell struct {
	// Date is the civil date represented by this position when HasDate is true.
	Date Date
	// HasDate is false only when an adjacent date crosses the supported year bounds.
	HasDate bool
	// Day is the numeric day of Date, even when its number is hidden.
	Day int
	// ShowDayNumber tells a renderer whether to display Day.
	ShowDayNumber bool
	// InMonth distinguishes requested-month dates from adjacent dates.
	InMonth bool
	// Weekday is Date's weekday when HasDate is true.
	Weekday time.Weekday
	// Weekend reports whether Weekday is Saturday or Sunday.
	Weekend bool
	// Events contains a copy of labels attached to Date.
	Events []string
	// ISOWeek is the ISO-8601 week number containing Date.
	ISOWeek int
	// ISOYear is the ISO-8601 week-year containing Date.
	ISOYear int
}

// Week is one complete row of seven calendar cells.
//
// ISOWeek and ISOYear describe the ISO week containing the Monday of this
// row, independent of the selected display week start.
type Week struct {
	// Cells are ordered from MonthGrid.WeekStart and always contain seven positions.
	Cells [7]Cell
	// ISOWeek is the ISO week containing Monday in this display row.
	ISOWeek int
	// ISOYear is the ISO week-year containing Monday in this display row.
	ISOYear int
}

// MonthGrid is a deterministic five- or six-row grid for a Gregorian month.
type MonthGrid struct {
	// Year is the Gregorian year represented by the grid.
	Year int
	// Month is the Gregorian month represented by the grid.
	Month time.Month
	// WeekStart is the first weekday in each row.
	WeekStart time.Weekday
	// ShowWeekNumbers copies the corresponding Options setting.
	ShowWeekNumbers bool
	// WeekdayLabels are ordered to match each row's cells.
	WeekdayLabels [7]string
	// Weeks contains five or six complete rows in display order.
	Weeks []Week
}

// NewMonthGrid constructs the complete week rows that contain a month.
//
// FixedRows may be 0 (minimum rows), 5, or 6. An error is returned if a
// requested fixed height cannot contain all dates, if labels are malformed,
// or if the year/month is outside the supported civil-calendar range.
func NewMonthGrid(year int, month time.Month, options Options) (MonthGrid, error) {
	if year < minYear || year > maxYear || month < time.January || month > time.December {
		return MonthGrid{}, fmt.Errorf("calendarboundary: invalid month %d/%d", year, month)
	}
	if options.WeekStart < time.Sunday || options.WeekStart > time.Saturday {
		return MonthGrid{}, fmt.Errorf("calendarboundary: invalid week start %d", options.WeekStart)
	}
	if options.FixedRows != 0 && options.FixedRows != 5 && options.FixedRows != 6 {
		return MonthGrid{}, fmt.Errorf("calendarboundary: fixed rows must be 0, 5, or 6")
	}
	if len(options.WeekdayLabels) != 0 && len(options.WeekdayLabels) != 7 {
		return MonthGrid{}, fmt.Errorf("calendarboundary: weekday labels must contain seven values")
	}
	first := Date{Year: year, Month: month, Day: 1}
	leading := (int(first.Weekday()) - int(options.WeekStart) + 7) % 7
	dayCount := DaysInMonth(year, month)
	requiredRows := max(5, (leading+dayCount+6)/7)
	rows := requiredRows
	if options.FixedRows != 0 {
		if requiredRows > options.FixedRows {
			return MonthGrid{}, fmt.Errorf("calendarboundary: month requires %d rows, cannot fit in %d", requiredRows, options.FixedRows)
		}
		rows = options.FixedRows
	}
	grid := MonthGrid{Year: year, Month: month, WeekStart: options.WeekStart, ShowWeekNumbers: options.WeekNumbers}
	for column := range grid.WeekdayLabels {
		weekday := (int(options.WeekStart) + column) % 7
		if len(options.WeekdayLabels) != 0 {
			grid.WeekdayLabels[column] = options.WeekdayLabels[weekday]
		} else {
			grid.WeekdayLabels[column] = time.Weekday(weekday).String()[:3]
		}
	}
	grid.Weeks = make([]Week, rows)
	for row := 0; row < rows; row++ {
		week := Week{}
		for column := 0; column < 7; column++ {
			slot := row*7 + column - leading + 1
			date := Date{Year: year, Month: month, Day: slot}
			inMonth := slot >= 1 && slot <= dayCount
			if !inMonth {
				date, _ = civilDateAtDayOffset(first, slot-1)
			}
			cell := Cell{Date: date, HasDate: date.Valid(), InMonth: inMonth}
			if cell.HasDate {
				cell.Day = date.Day
				cell.Weekday = date.Weekday()
				cell.Weekend = cell.Weekday == time.Saturday || cell.Weekday == time.Sunday
				cell.ISOYear, cell.ISOWeek = ISOWeek(date)
				cell.ShowDayNumber = !options.Blank && (inMonth || options.ShowOutsideDays)
				if labels, ok := options.Events[date]; ok && (inMonth || options.ShowOutsideDays) {
					cell.Events = append([]string(nil), labels...)
				}
			}
			week.Cells[column] = cell
		}
		mondayColumn := (int(time.Monday) - int(options.WeekStart) + 7) % 7
		mondayCell := week.Cells[mondayColumn]
		if mondayCell.HasDate {
			week.ISOYear, week.ISOWeek = mondayCell.ISOYear, mondayCell.ISOWeek
		}
		grid.Weeks[row] = week
	}
	return grid, nil
}

// NewYearGrid builds one month grid for every month in the requested year.
func NewYearGrid(year int, options Options) ([12]MonthGrid, error) {
	var grids [12]MonthGrid
	if year < minYear || year > maxYear {
		return grids, fmt.Errorf("calendarboundary: invalid year %d", year)
	}
	for index := range grids {
		grid, err := NewMonthGrid(year, time.Month(index+1), options)
		if err != nil {
			return [12]MonthGrid{}, err
		}
		grids[index] = grid
	}
	return grids, nil
}

func civilDateAtDayOffset(start Date, offset int) (Date, bool) {
	year, month, day := civilFromDays(daysFromCivil(start.Year, int(start.Month), start.Day) + offset)
	date := Date{Year: year, Month: time.Month(month), Day: day}
	return date, date.Valid()
}
