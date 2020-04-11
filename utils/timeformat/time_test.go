package timeformat

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDateTimeFormat(t *testing.T) {
	s := "2007-01-02 15:04:05"
	tm, err := ParseLongDate(s)
	require.Nil(t, err)
	require.Equal(t, s, LongDateFormat(tm))
}

func TestShortDateTime(t *testing.T) {
	s := "2007-01-02"
	tm, err := ParseShortDate(s)
	require.Nil(t, err)
	require.Equal(t, s, ShortDateFormat(tm))
}