package calendarboundary

import (
	"fmt"
	"time"
)

func ExampleNewMonthGrid() {
	grid, _ := NewMonthGrid(2028, time.February, Options{WeekStart: time.Monday, ShowOutsideDays: true})
	fmt.Printf("%d rows, %d days, last day %s\n", len(grid.Weeks), DaysInMonth(grid.Year, grid.Month), grid.Weeks[4].Cells[0].Date)
	// Output: 5 rows, 29 days, last day 2028-02-28
}
