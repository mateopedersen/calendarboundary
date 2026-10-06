// Package calendarboundary computes civil-date and calendar-grid boundaries
// without depending on local time zones or a rendering framework.
//
// Dates use the proleptic Gregorian calendar and years 1 through 9999. Month
// grids have deterministic complete weeks, configurable week starts, optional
// outside-month numbers and event labels, and ISO week metadata. The package
// is intended to provide stable data for printers, planning views, and APIs.
package calendarboundary
