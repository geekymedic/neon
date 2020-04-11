package timeformat

import (
	"time"

	"github.com/geekymedic/neon/pkg/types"
)

func LongDateFormat(t time.Time) string {
	return t.Format(types.LongDateType)
}

func ShortDateFormat(t time.Time) string {
	return t.Format(types.ShortDateType)
}

func ParseLongDate(s string) (t time.Time, err error) {
	return time.ParseInLocation(types.LongDateType, s, time.Local)
}

func ParseShortDate(s string) (t time.Time, err error) {
	return time.ParseInLocation(types.ShortDateType, s, time.Local)
}
