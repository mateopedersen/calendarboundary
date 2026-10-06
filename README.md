# CalendarBoundary

CalendarBoundary is a small Go library for Gregorian civil dates, month-grid construction, and ISO week metadata. It returns calendar data rather than HTML or UI widgets, so applications can render the same computed boundaries in a terminal, PDF, API, or planning interface.

## Features

- Validated civil dates and timezone-free date arithmetic for years 1–9999.
- Gregorian leap-year and month-length rules, including century and 400-year cases.
- Five- or six-row month grids with configurable Sunday-through-Saturday starts.
- Per-cell weekday, weekend, date, ISO week-year, and event metadata.
- Optional adjacent-month labels, blank planning mode, custom weekday labels, and fixed grid heights.
- A 12-month year grid with no third-party dependencies.

## Why calendar boundaries are hard

Month lengths vary, century years do not always leap, and a month can begin on any weekday. A display row can also straddle an ISO week-year boundary: for example, January 1 can belong to the previous ISO week-year. Week-start conventions change column positions but must not change ISO metadata. A robust grid needs to make each of these decisions explicit and stable.

## Install

```sh
go get github.com/mateopedersen/calendarboundary@v0.1.0
```

## Quick start

```go
package main

import (
	"fmt"
	"time"

	"github.com/mateopedersen/calendarboundary"
)

func main() {
	grid, err := calendarboundary.NewMonthGrid(2028, time.February, calendarboundary.Options{
		WeekStart:       time.Monday,
		ShowOutsideDays: true,
		WeekNumbers:     true,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d days in %d rows\n", calendarboundary.DaysInMonth(grid.Year, grid.Month), len(grid.Weeks))
	fmt.Printf("last date: %s\n", grid.Weeks[4].Cells[1].Date)
}
```

## Grid model

`NewMonthGrid` creates five complete rows by default and uses six when the month needs them. `Options.FixedRows` can request a five- or six-row shape; a request that cannot contain the month returns an error. Each `Week` has seven cells ordered from `Options.WeekStart`.

Every in-range `Cell` contains an actual `Date`, weekday, weekend flag, ISO week number and year, and a copied slice of event labels. `ShowDayNumber` lets renderers preserve the position while hiding adjacent-month or blank-mode numbers. Dates outside the documented 1–9999 range are marked with `HasDate == false` at the calendar's outer edge.

## Week starts and ISO week-years

`WeekStart` uses Go's `time.Weekday` values. Choose `time.Monday` or `time.Sunday` for common conventions; any weekday is supported. `Cell.ISOWeek` and `Cell.ISOYear` always use ISO-8601's Monday-first week definition. Row-level metadata uses the ISO week containing that row's Monday.

```go
date, _ := calendarboundary.NewDate(2021, time.January, 1)
isoYear, week := date.ISOWeek()
fmt.Printf("%d-W%02d\n", isoYear, week) // 2020-W53
```

## Leap years and date arithmetic

`IsLeapYear` follows the Gregorian 4/100/400 rule. `DaysInMonth` returns zero for an invalid year or month. `NewDate` validates inputs, and `Date.AddDays` returns an error if arithmetic would leave years 1–9999.

## Year grids and events

`NewYearGrid(year, options)` returns a fixed array of twelve month grids in January-to-December order. Events are supplied as `map[Date][]string`; labels are copied into matching visible cells, so later mutation of the input slice does not change the cell.

## Testing

```sh
go test ./...
go vet ./...
go test -race ./...
```

The tests include leap centuries, leap days, December-to-January arithmetic, ISO year crossovers, Monday/Sunday grid starts, five/six-row behavior, blank cells, and year overviews.

## Calendar resources

This library focuses on calendar computation. Human-facing printable and planning calendar resources are available through [BetaCalendars](https://www.betacalendars.com/).

## License

MIT. See [LICENSE](LICENSE).
