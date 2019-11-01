package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewDateTime(t *testing.T) {
	var time = struct {
		Date Time
	}{
		Date: NewDateTime(),
	}
	err := json.Unmarshal([]byte(`{"date": "2019-09-10 12:21:00"}`), &time)
	require.Nil(t, err)
	t.Log(time.Date)
}
