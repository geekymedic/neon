package v8

import (
	"testing"
)

func TestIsSpecialName(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		if ok := IsNickName(nil,  "日本語abc_0123456789", ""); !ok {
			t.Fatal("expect true, actual false")
		}
	})
	t.Run("fail", func(t *testing.T) {
		if ok := IsNickName(nil, "#", ""); ok {
			t.Fatal("expect false, actual true")
		}

		if ok := IsNickName(nil, "-", ""); ok {
			t.Fatal("expect false, actual true")
		}

		if ok := IsNickName(nil, "?", ""); ok {
			t.Fatal("expect false, actual true")
		}
	})
}
