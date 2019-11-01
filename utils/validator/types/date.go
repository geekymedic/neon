package types

import (
	"fmt"
	"strings"
	t "time"
)

func NewDateTime() Time {
	return Time{
		Time:   t.Time{},
		format: "2006-01-02 15:04:05",
	}
}

func NewDateTimeWithTime(tm t.Time) *Time {
	return &Time{
		Time:   tm,
		format: "2006-01-02 15:04:05",
	}
}

type Time struct {
	Time   t.Time
	format string
}

func (time *Time) UnmarshalJSON(b []byte) error {
	var tm, err = t.ParseInLocation(time.format, strings.Trim(fmt.Sprintf("%s", b), "\""), t.Local)
	if err != nil {
		return err
	}
	time.Time = tm
	return nil
}
