// Description: This package provides db management for the application.
// Developer: Aleksei Grigorev <https://github.com/AlekseiGrigorev>, <aleksvgrig@gmail.com>
// Copyright (c) 2025 Aleksei Grigorev
package appdb

import (
	"strings"
	"time"
)

// QueryProcessorBetween is a struct that implements the QueryProcessorInterface interface
// for processing SQL queries with BETWEEN clauses for date ranges.
type QueryProcessorBetween struct {
	Query        string
	Start        string
	End          string
	Step         string
	currentStart string
}

// Return the type name for a query processor.
func (q *QueryProcessorBetween) GetType() string {
	return QUERY_TYPE_BETWEEN
}

// InitQuery resets the query processor to its initial state by setting the
// values of the Start, End, and Step fields to empty strings and the
// currentStart field to an empty string.
func (q *QueryProcessorBetween) InitQuery() QueryProcessorInterface {
	q.Start = ""
	q.End = ""
	q.Step = ""
	q.currentStart = ""
	return q
}

// setValue sets the value of the specified key in the query processor.
func (q *QueryProcessorBetween) SetValue(key string, value any) QueryProcessorInterface {
	switch strings.ToLower(key) {
	case "start":
		q.Start = value.(string)
	case "end":
		q.End = value.(string)
	case "step":
		q.Step = value.(string)
	}
	return q
}

// ProcessQuery implements the QueryProcessorInterface and returns the query string
// with the {{start}} and {{end}} placeholders replaced by the current start and end
// dates, respectively. The current start date is incremented by the step
// duration after each call, facilitating paging through results. If the current end
// date is after the end date specified in the query processor, the current end
// date is set to the end date. The function returns the resulting query string
// with the placeholders replaced.
func (q *QueryProcessorBetween) ProcessQuery() string {
	trimmedQuery := strings.TrimRight(q.Query, " \t\n\r;")
	start, end := q.getBetweenDates()
	return strings.ReplaceAll(strings.ReplaceAll(trimmedQuery, "{{start}}", start), "{{end}}", end)
}

// getBetweenDates returns the start and end dates for the BETWEEN clause of the query string.
// The function increments the current start date by the step duration after each call, facilitating
// paging through results. If the current end date is after the end date specified in the query
// processor, the current end date is set to the end date. The function returns the resulting start and
// end dates as strings in the format specified by DATE_TIME_LAYOUT.
func (q *QueryProcessorBetween) getBetweenDates() (string, string) {
	if q.currentStart == "" {
		q.currentStart = q.Start
	}
	currentStart, err := time.Parse(DATE_TIME_LAYOUT, q.currentStart)
	if err != nil {
		return "", ""
	}
	duration, err := time.ParseDuration(q.Step)
	if err != nil {
		return "", ""
	}
	currentEnd := currentStart.Add(duration)
	end, err := time.Parse(DATE_TIME_LAYOUT, q.End)
	if err != nil {
		return "", ""
	}
	if currentEnd.After(end) {
		currentEnd = end
	}
	if !currentStart.After(end) {
		q.currentStart = currentStart.Add(duration).Format(DATE_TIME_LAYOUT)
	}
	return currentStart.Format(DATE_TIME_LAYOUT), currentEnd.Format(DATE_TIME_LAYOUT)
}
