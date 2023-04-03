package blevExp

import (
	"github.com/blevesearch/bleve/v2/analysis"
	"time"
)

const rfc3339NoTimezone = "2006-01-02T15:04:05"
const rfc3339NoTimezoneNoT = "2006-01-02 15:04:05"
const rfc3339NoTime = "2006-01-02"

var layouts = []string{
	time.Layout,
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339,
	time.RFC3339Nano,
	rfc3339NoTimezone,
	rfc3339NoTimezoneNoT,
	rfc3339NoTime,
	"02 Jan 2006",
	"2006-01-02",
	"2006-01-02T15:04:05Z07:00",
	"Jan 02, 2006",
	"January 02, 2006",
	"January 02th, 2006",
	"Mon, 02 Jan 2006 03:04:05 +0200",
	"Mon, 02 Jan 2006 15:04:05 -0300",
	"Mon, 02 Jan 2006 15:04:05 MST",
}

func ParseDateTime(input string) (time.Time, error) {
	for _, layout := range layouts {
		rv, err := time.Parse(layout, input)
		if err == nil {
			return rv, nil
		}
	}
	return time.Time{}, analysis.ErrInvalidDateTime
}
