package validator

import (
	"reflect"
	"testing"
	"time"
)

func validate(field *Field, fn ...OpFn) error {
	var val = NewValidator()
	for _, f := range fn {
		err := val.Validate(field, f).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

func TestEq(t *testing.T) {
	SetValidatorVersion(v8)
	t.Run("v8, not eq", func(t *testing.T) {
		// string
		{
			err := validate(&Field{Tag: "address", Value: "BeiJin"}, Eq("ShangHai"))
			if err == nil {
				t.Fatalf("expect: %v, actual: nil", err)
			}
		}

		// int
		{
			err := validate(&Field{Tag: "age", Value: 100}, Eq("101"))
			if err == nil {
				t.Fatalf("expect: %v, actual: nil", err)
			}
		}

		// float
		{
			err := validate(&Field{Tag: "price", Value: 100.203}, Eq("100.205"))
			if err == nil {
				t.Fatalf("expect: %v, actual: nil", err)
			}
		}

		{
			err := validate(&Field{Tag: "children", Value: []string{"zhanshan", "lisi"}}, Eq("3"))
			if err == nil {
				t.Fatalf("expect: %v, actual: nil", err)
			}
		}
	})
	t.Run("v8, eq", func(t *testing.T) {
		{
			err := validate(&Field{Tag: "address", Value: "BeiJin"}, Eq("BeiJin"))
			if err != nil {
				t.Fatalf("expect: nil, actual: %v", err)
			}
		}

		{
			err := validate(&Field{Tag: "age", Value: 100}, Eq("100"))
			if err != nil {
				t.Fatalf("expect: nil, actual: %v", err)
			}
		}

		// float
		{
			err := validate(&Field{Tag: "price", Value: 100.203}, Eq("100.203"))
			if err != nil {
				t.Fatalf("expect: nil, actual: %v", err)
			}
		}

		// slice
		{
			err := validate(&Field{Tag: "children", Value: []string{"zhanshan", "lisi"}}, Eq("2"))
			if err != nil {
				t.Fatalf("expect: nil, actual: %v", err)
			}
		}
	})
}
func TestMax(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		tm := time.Now().Add(time.Second)
		args := map[string]interface{}{
			"0":         []int{1},                                     // len
			"2":         map[string]string{"a": "", "b": "", "c": ""}, // len
			"3":         4,                                            // value
			"4.1":       4.2,                                          // value
			"5":         [6]int{},                                     // len
			"6":         [7]int{},
			tm.String(): tm, // Precision is relatively fine for time, so add seconds for test
		}

		for key, value := range args {
			err := validate(&Field{"test", value}, Max(key))
			if err == nil {
				t.Fatalf("expect: %v, actual: nil", err)
			}
		}
	})

	t.Run("ok", func(t *testing.T) {
		tm := time.Now()
		args := map[string]interface{}{
			"2":         []int{1, 2},
			"3":         map[int]int{1: 1},
			"4":         3,
			"5":         4.9,
			"6":         [5]int{},
			tm.String(): tm, // Precision is relatively fine for time, so add seconds for test
		}
		for key, value := range args {
			err := validate(&Field{"test", value}, Max(key))
			if err != nil {
				t.Fatalf("expect: nil, actual: %v", err)
			}
		}
	})
}

func TestMin(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		tm := time.Now().Add(time.Second)
		args := map[string]interface{}{
			"0":         []int{1},                                     // len
			"2":         map[string]string{"a": "", "b": "", "c": ""}, // len
			"3":         4,                                            // value
			"4.1":       4.2,                                          // value
			"5":         [6]int{},                                     // len
			"6":         [7]int{},
			tm.String(): tm, // Precision is relatively fine for time, so add seconds for test
		}

		for key, value := range args {
			err := validate(&Field{"test", value}, Min(key))
			if err != nil {
				t.Fatalf("expect: nil, actual: %v", err)
			}
		}
	})

	t.Run("ok", func(t *testing.T) {
		tm := time.Now().Add(time.Second)
		args := map[string]interface{}{
			"2":         []int{1, 2},
			"3":         map[int]int{1: 1, 2: 2, 3: 3, 4: 4, 5: 5},
			"4":         5,
			"5":         5.9,
			"6":         [6]int{},
			tm.String(): tm, // Precision is relatively fine for time, so add seconds for test
		}
		for key, value := range args {
			err := validate(&Field{"test", value}, Min(key))
			if err != nil {
				t.Fatalf("expect: nil, actual: %v", err)
			}
		}
	})
}

func TestRange(t *testing.T) {
	t.Run("fail, [min, max]", func(t *testing.T) {
		err := validate(&Field{"", 21}, Min("30"), Max("39"))
		if err == nil {
			t.Fatalf("expect: %v, actual: nil", err)
		}

		err = validate(&Field{"", 41}, Min("30"), Max("39"))
		if err == nil {
			t.Fatalf("expect: %v, actual: nil", err)
		}
	})

	t.Run("ok, [min, max]", func(t *testing.T) {
		err := validate(&Field{"", 21}, Min("20"), Max("39"))
		if err != nil {
			t.Fatalf("expect: nil, actual: %v", err)
		}

		err = validate(&Field{"", 21}, Min("39"), Max("39"))
		if err == nil {
			t.Fatalf("expect: nil, actual: %v", err)
		}
	})
}

func TestGt(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		err := validate(&Field{"Age", 10}, Gt("20"))
		if err == nil {
			t.Fatalf("expect: %v, actual: nil", err)
		}
	})

	t.Run("ok", func(t *testing.T) {
		err := validate(&Field{"Age", 10}, Gt("3"))
		if err != nil {
			t.Fatalf("expect: nil, actual: %v", err)
		}
	})
}

func TestGte(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		tm := time.Now()
		args := map[string]interface{}{
			"2":         []int{1},
			"3":         map[int]int{1: 1},
			"4":         3,
			"5":         4.9,
			"6":         [5]int{},
			tm.String(): tm, // Precision is relatively fine for time, so add seconds for test
		}
		for key, value := range args {
			err := validate(&Field{"test", value}, Gte(key))
			if err == nil {
				t.Fatalf("expect: %v, actual: nil", err)
			}
		}
	})

	t.Run("ok", func(t *testing.T) {
		tm := time.Now().Add(time.Second)
		args := map[string]interface{}{
			"1":         []int{1},                            // len
			"2":         map[string]string{"a": "", "b": ""}, // len
			"3":         3,                                   // value
			"4.1":       4.1,                                 // value
			"5":         [5]int{},                            // len
			"6":         [7]int{},
			tm.String(): tm, // Precision is relatively fine for time, so add seconds for test
		}
		for key, value := range args {
			err := validate(&Field{"test", value}, Gte(key))
			if err != nil {
				t.Fatalf("expect: nil, actual: %v", err)
			}
		}
	})
}

func TestLt(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		tm := time.Now().Add(time.Second)
		args := map[string]interface{}{
			"1":         []int{1},                            // len
			"2":         map[string]string{"a": "", "b": ""}, // len
			"3":         3,                                   // value
			"4.1":       4.1,                                 // value
			"5":         [5]int{},                            // len
			"6":         [7]int{},
			tm.String(): tm, // Precision is relatively fine for time, so add seconds for test
		}
		for key, value := range args {
			err := validate(&Field{"test", value}, Lt(key))
			if err == nil {
				t.Fatalf("expect: %v, actual: nil", err)
			}
		}
	})

	t.Run("ok", func(t *testing.T) {
		tm := time.Now()
		args := map[string]interface{}{
			"1":         []int{},                    // len
			"2":         map[string]string{"a": ""}, // len
			"3":         2,                          // value
			"4.1":       4.0,                        // value
			"5":         [4]int{},                   // len
			"6":         [5]int{},
			tm.String(): tm, // Precision is relatively fine for time, so add seconds for test
		}
		for key, value := range args {
			err := validate(&Field{"test", value}, Lt(key))
			if err != nil {
				t.Fatalf("expect: nil, actual: %v", err)
			}
		}
	})
}

func TestLte(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		tm := time.Now().Add(time.Second)
		args := map[string]interface{}{
			"1":         []int{1, 2},                                  // len
			"2":         map[string]string{"a": "", "b": "", "c": ""}, // len
			"3":         4,                                            // value
			"4.1":       5.1,                                          // value
			"5":         [6]int{},                                     // len
			"6":         [8]int{},
			tm.String(): tm, // Precision is relatively fine for time, so add seconds for test
		}
		for key, value := range args {
			err := validate(&Field{"test", value}, Lte(key))
			if err == nil {
				t.Fatalf("expect: %v, actual: nil", err)
			}
		}
	})

	t.Run("ok", func(t *testing.T) {
		tm := time.Now()
		args := map[string]interface{}{
			"1":         []int{1},                   // len
			"2":         map[string]string{"a": ""}, // len
			"3":         2,                          // value
			"4.1":       4.0,                        // value
			"5":         [4]int{},                   // len
			"6":         [5]int{},
			tm.String(): tm, // Precision is relatively fine for time, so add seconds for test
		}
		for key, value := range args {
			err := validate(&Field{"test", value}, Lte(key))
			if err != nil {
				t.Fatalf("expect: nil, actual: %v", err)
			}
		}
	})
}

func TestContains(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		args := map[string]interface{}{
			"c": "abd",
		}
		for key, value := range args {
			err := validate(&Field{"test", value}, Contains(key))
			if err == nil {
				t.Fatalf("expect: %v, actual: nil", err)
			}
		}
	})

	t.Run("ok", func(t *testing.T) {
		args := map[string]interface{}{
			"bd": "abd",
		}
		for key, value := range args {
			err := validate(&Field{"test", value}, Contains(key))
			if err != nil {
				t.Fatalf("expect: nil, actual: %v", err)
			}
		}
	})
}

func TestContainsAny(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		args := map[string]interface{}{
			"c": "abd",
		}
		for key, value := range args {
			err := validate(&Field{"test", value}, ContainsAny(key))
			if err == nil {
				t.Fatalf("expect: %v, actual: nil", err)
			}
		}
	})

	t.Run("ok", func(t *testing.T) {
		args := map[string]interface{}{
			"ad":   "abd",
			"a":    "abd",
			"abdc": "abd",
		}
		for key, value := range args {
			err := validate(&Field{"test", value}, ContainsAny(key))
			if err != nil {
				t.Fatalf("expect: nil, actual: %v", err)
			}
		}
	})
}

func TestContainsRune(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		args := map[string]interface{}{
			"中": "日本",
		}
		for key, value := range args {
			err := validate(&Field{"test", value}, ContainsRune(key))
			if err == nil {
				t.Fatalf("expect: %v, actual: nil", err)
			}
		}
	})

	t.Run("ok", func(t *testing.T) {
		args := map[string]interface{}{
			"北":  "北京",
			"中国": "中国好!",
		}
		for key, value := range args {
			err := validate(&Field{"test", value}, ContainsAny(key))
			if err != nil {
				t.Fatalf("expect: nil, actual: %v", err)
			}
		}
	})
}

func TestExcludes(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		args := []struct {
			value string
			param string
		}{
			{"abcs", "abc"},
			{"dddd", "dd"},
		}
		for _, arg := range args {
			err := validate(&Field{"test", arg.value}, Excludes(arg.param))
			if err == nil {
				t.Fatalf("expect: not nil, actual: nil")
			}
		}
	})

	t.Run("ok", func(t *testing.T) {
		args := []struct {
			value string
			param string
		}{{"accc", "acd"}, {"adad", "dd"}}
		for _, arg := range args {
			err := validate(&Field{"test", arg.value}, Excludes(arg.param))
			if err != nil {
				t.Fatalf("expect: nil, actual: %v", err)
			}
		}
	})
}

func TestExcludesAll(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		args := map[string]interface{}{
			"abc": "abcs",
			"dd":  "ddd",
		}
		for key, value := range args {
			err := validate(&Field{"test", value}, ExcludesAll(key))
			if err == nil {
				t.Fatalf("expect: %v, actual: nil", err)
			}
		}
	})

	t.Run("ok", func(t *testing.T) {
		args := map[string]interface{}{
			"abc": "defg",
			"dd":  "aaaa",
		}
		for key, value := range args {
			err := validate(&Field{"test", value}, ExcludesAll(key))
			if err != nil {
				t.Fatalf("expect: nil, actual: %v", err)
			}
		}
	})
}

func TestExcludesRune(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		args := map[string]interface{}{
			"中": "中国",
		}
		for key, value := range args {
			err := validate(&Field{"test", value}, ExcludesRune(key))
			if err == nil {
				t.Fatalf("expect: %v, actual: nil", err)
			}
		}
	})

	t.Run("ok", func(t *testing.T) {
		args := map[string]interface{}{
			"北":  "上海",
			"美国": "中国好!",
		}
		for key, value := range args {
			err := validate(&Field{"test", value}, ExcludesRune(key))
			if err != nil {
				t.Fatalf("expect: nil, actual: %v", err)
			}
		}
	})
}

func TestIsMAC(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		err := validate(&Field{"Mac", "44-45-53-54-00-00-00"}, IsMAC())
		if err == nil {
			t.Fatalf("expect: %v, actual: nil", err)
		}
	})

	t.Run("ok", func(t *testing.T) {
		err := validate(&Field{"Mac", "44-45-53-54-00-00"}, IsMAC())
		if err != nil {
			t.Fatalf("expect: nil, actual: %v", err)
		}
	})
}

func TestIsBase64(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		err := validate(&Field{"", "aGVsbG93b3Jk-"}, IsBase64())
		if err == nil {
			t.Fatalf("epxect: %v, actual: nil", err)
		}
	})

	t.Run("ok", func(t *testing.T) {
		err := validate(&Field{"", "aGVsbG93b3Jk"}, IsBase64())
		if err != nil {
			t.Fatalf("expect: nil, actual: %v", err)
		}
	})
}

func TestIsIP(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		err := validate(&Field{"IP", "127.0.0.0.0"}, IsIP())
		if err == nil {
			t.Fatalf("expect: %v, actual: nil", err)
		}
	})

	t.Run("ok", func(t *testing.T) {
		err := validate(&Field{"IP", "127.0.0.1"}, IsIP())
		if err != nil {
			t.Fatalf("expect: nil, actual: %v", err)
		}
	})
}

func TestIsIPv4(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		err := validate(&Field{"IP", "127.0.0.0.0"}, IsIPv4())
		if err == nil {
			t.Fatalf("expect: %v, actual: nil", err)
		}
	})

	t.Run("ok", func(t *testing.T) {
		err := validate(&Field{"IP", "183.14.132.105"}, IsIPv4())
		if err != nil {
			t.Fatalf("expect: nil, actual: %v", err)
		}
	})
}

func TestIsIPv6(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		err := validate(&Field{"IP", "183.14.132.105"}, IsIPv6())
		if err == nil {
			t.Fatalf("expect: %v, actual: nil", err)
		}
	})

	t.Run("ok", func(t *testing.T) {
		err := validate(&Field{"IP", "2001:da8:8000:1::81"}, IsIPv6())
		if err != nil {
			t.Fatalf("expect: nil, actual: %v", err)
		}
	})
}

func TestIsAlpha(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		err := validate(&Field{"IsAlpha", "183.14.132.105"}, IsAlpha())
		if err == nil {
			t.Fatalf("expect: %v, actual: nil", err)
		}
	})

	t.Run("ok", func(t *testing.T) {
		err := validate(&Field{"IsAlpha", "aljalsjldjlakjsl"}, IsAlpha())
		if err != nil {
			t.Fatalf("expect: nil, actual: %v", err)
		}
	})
}

func TestIsAlphanum(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		err := validate(&Field{"IsAlpha", "183.14.132.105"}, IsAlphanum())
		if err == nil {
			t.Fatalf("expect: %v, actual: nil", err)
		}
	})

	t.Run("ok", func(t *testing.T) {
		err := validate(&Field{"IsAlpha", "112aljalsjldjl123akjsl"}, IsAlphanum())
		if err != nil {
			t.Fatalf("expect: nil, actual: %v", err)
		}
	})
}

func TestIsNumber(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		err := validate(&Field{"IsNumber", "183d"}, IsNumber())
		if err == nil {
			t.Fatalf("expect: %v, actual: nil", err)
		}
	})

	t.Run("ok", func(t *testing.T) {
		err := validate(&Field{"IsNumber", "00344"}, IsNumber())
		if err != nil {
			t.Fatalf("expect: nil, actual: %v", err)
		}
	})
}

func TestIsNumeric(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		err := validate(&Field{"IsNumeric", "183d"}, IsNumeric())
		if err == nil {
			t.Fatalf("expect: %v, actual: nil", err)
		}
	})

	t.Run("ok", func(t *testing.T) {
		arg := []struct {
			param  string
			expect error
		}{
			{"+01200344", nil},
			{"+0.1200344", nil},
			{"-01200344", nil},
			{"-0.1200344", nil},
		}

		for _, arg := range arg {
			err := validate(&Field{"IsNumeric", arg.param}, IsNumeric())
			if err != nil {
				t.Fatalf("expect: %v, actual: %v", arg.expect, err)
			}
		}
	})
}

func TestIsIP4AddrResolvable(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		err := validate(&Field{"IP", "2001:da8:8000:1::81"}, IsIP4AddrResolvable())
		if err == nil {
			t.Fatalf("expect: %v, actual: nil", err)
		}
	})

	t.Run("ok", func(t *testing.T) {
		err := validate(&Field{"IP", "127.0.0.1"}, IsIP4AddrResolvable())
		if err != nil {
			t.Fatalf("expect: nil, actual: %v", err)
		}
	})
}

func TestIsIP6AddrResolvable(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		err := validate(&Field{"IP", "103.123.3.1"}, IsIP6AddrResolvable())
		if err == nil {
			t.Fatalf("expect: %v, actual: nil", err)
		}
	})

	t.Run("ok", func(t *testing.T) {
		err := validate(&Field{"IP", "2001:da8:8000:1::81"}, IsIP6AddrResolvable())
		if err != nil {
			t.Fatalf("expect: nil, actual: %v", err)
		}
	})
}

func TestIsIPAddrResolvable(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		err := validate(&Field{"IP", "103.123.3.1.0"}, IsIPAddrResolvable())
		if err == nil {
			t.Fatalf("expect: %v, actual: nil", err)
		}
	})

	t.Run("ok", func(t *testing.T) {
		args := []string{
			"127.0.0.1",
			"::1",
			"2001:da8:8000:1::81",
		}
		for _, arg := range args {
			err := validate(&Field{"IP", arg}, IsIPAddrResolvable())
			if err != nil {
				t.Fatalf("expect: nil, actual: %v", err)
			}
		}
	})
}

func TestIsTCPAddrResolvable(t *testing.T) {
	t.Run("tcp", func(t *testing.T) {
		for _, value := range []string{"256.0.0.0:1", "[::1]", ":80"} {
			err := validate(&Field{"tcp", value}, IsTCPAddrResolvable())
			if err == nil {
				t.Fatalf("expect: %v, actual: nil", err)
			}
		}

		for _, value := range []string{"[::1]:80", "127.0.0.1:80"} {
			err := validate(&Field{"tcp", value}, IsTCPAddrResolvable())
			if err != nil {
				t.Fatalf("expect: nil, actual: %v", err)
			}
		}
	})
	t.Run("tcp4", func(t *testing.T) {
		for _, value := range []string{"[::1]:80", ":80"} {
			err := validate(&Field{"tcp", value}, IsTCP4AddrResolvable())
			if err == nil {
				t.Fatalf("expect: %v, actual: nil", err)
			}
		}

		for _, value := range []string{"127.0.0.1:80"} {
			err := validate(&Field{"tcp", value}, IsTCP4AddrResolvable())
			if err != nil {
				t.Fatalf("expect: nil, actual: %v", err)
			}
		}
	})
	t.Run("tcp6", func(t *testing.T) {
		for _, value := range []string{"[::1]", ":80"} {
			err := validate(&Field{"tcp", value}, IsTCP6AddrResolvable())
			if err == nil {
				t.Fatalf("expect: %v, actual: nil", err)
			}
		}

		for _, value := range []string{"[::1]:80"} {
			err := validate(&Field{"tcp", value}, IsTCP6AddrResolvable())
			if err != nil {
				t.Fatalf("expect: nil, actual: %v", err)
			}
		}
	})
}

func TestIsUDPAddrResolvable(t *testing.T) {
	t.Run("udp", func(t *testing.T) {
		for _, value := range []string{":80", "[::1]", "256.0.0.0:1"} {
			err := validate(&Field{"tcp", value}, IsUDPAddrResolvable())
			if err == nil {
				t.Fatalf("expect: %v, actual: nil", err)
			}
		}

		for _, value := range []string{"127.0.0.1:80", "[::1]:80"} {
			err := validate(&Field{"tcp", value}, IsUDPAddrResolvable())
			if err != nil {
				t.Fatalf("expect: nil, actual: %v", err)
			}
		}
	})
	t.Run("udp4", func(t *testing.T) {
		for _, value := range []string{"[::1]:80"} {
			err := validate(&Field{"tcp", value}, IsUDP4AddrResolvable())
			if err == nil {
				t.Fatalf("expect: %v, actual: nil", err)
			}
		}

		for _, value := range []string{"127.0.0.1:80"} {
			err := validate(&Field{"tcp", value}, IsUDP4AddrResolvable())
			if err != nil {
				t.Fatalf("expect: nil, actual: %v", err)
			}
		}
	})
	t.Run("tcp6", func(t *testing.T) {
		for _, value := range []string{"[::1]", ":80"} {
			err := validate(&Field{"tcp", value}, IsUDP6AddrResolvable())
			if err == nil {
				t.Fatalf("expect: %v, actual: nil", err)
			}
		}

		for _, value := range []string{"[::1]:80"} {
			err := validate(&Field{"tcp", value}, IsUDP6AddrResolvable())
			if err != nil {
				t.Fatalf("expect: nil, actual: %v", err)
			}
		}
	})
}

func TestIsUnixAddrResolvable(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		args := []struct {
			param  string
			expect error
		}{{"", nil}, {"v.sock", nil}}
		for _, arg := range args {
			err := validate(&Field{"IP", arg.param}, IsUnixAddrResolvable())
			if err != arg.expect {
				t.Fatalf("expect: %v, actual: %v", arg.expect, err)
			}
		}
	})
}

func TestIsLatitude(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		args := []struct {
			param  string
			expect error
		}{{"+90", nil}, {"-90", nil}}
		for _, arg := range args {
			err := validate(&Field{"Coordinates", arg.param}, IsLatitude())
			if err != nil {
				t.Fatalf("expect: %v, actual: %v", arg.expect, err)
			}
		}
	})
}

func TestIsLongitude(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		args := []struct {
			param  string
			expect error
		}{{"+130", nil}, {"-100", nil}}
		for _, arg := range args {
			err := validate(&Field{"Coordinates", arg.param}, IsLongitude())
			if err != nil {
				t.Fatalf("expect: %v, actual: %v", arg.expect, err)
			}
		}
	})
}

func TestIsEmail(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		err := validate(&Field{"Eamil", "hellowordgmail.com"}, IsEmail())
		if err == nil {
			t.Fatalf("expect: %v, actual: nil", err)
		}
	})

	t.Run("ok", func(t *testing.T) {
		err := validate(&Field{"Eamil", "helloword@gmail.com"}, IsEmail())
		if err != nil {
			t.Fatalf("expect: nil, actual: %v", err)
		}
	})
}

func TestLen(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		err := validate(&Field{"email", "helloword@gmail.com"}, Len("20"))
		if err == nil {
			t.Fatalf("expect: %v, actual: nil", err)
		}
	})

	t.Run("ok", func(t *testing.T) {
		err := validate(&Field{"email", "helloword@gmail.com"}, Len("19"))
		if err != nil {
			t.Fatalf("expect: nil, actual: %v", err)
		}
	})
}

func TestNumber(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		args := []string{
			"10.1",
			"a",
			"a10",
			"10.1a",
		}
		for _, arg := range args {
			err := validate(&Field{"Age", arg}, IsNumber())
			if err == nil {
				t.Fatalf("expect: %v, actual: nil", err)
			}
		}
	})

	t.Run("ok", func(t *testing.T) {
		args := []string{
			"10",
			"010",
			"19930102",
		}
		for _, arg := range args {
			err := validate(&Field{"Age", arg}, IsNumber())
			if err != nil {
				t.Fatalf("expect: nil, actual: %v", err)
			}
		}
	})
}

func TestIsTime(t *testing.T) {
	var err error
	if err = validate(&Field{"time", "2017-03-04 00:00:00"}, IsTime()); err != nil {
		t.Fatalf("expect: nil, got: %v", err)
	}

	if err = validate(&Field{"time", "2016-03-04 00:00:00"}, IsTime()); err != nil {
		t.Fatalf("expect: nil, got: %v", err)
	}

	if err = validate(&Field{"time", "2016-03-04 11:00:01"}, IsTime()); err != nil {
		t.Fatalf("expect: nil, got: %v", err)
	}
	if err = validate(&Field{"time", "2016-03-04 11:00:00"}, IsTime()); err != nil {
		t.Fatalf("expect: nil, got: %v", err)
	}
	if err = validate(&Field{"time", "2016-03-04 1:00:00"}, IsTime()); err != nil {
		t.Fatalf("expect: nil, got: %v", err)
	}
	if err = validate(&Field{"time", "2016-03-04 1:01:01"}, IsTime()); err != nil {
		t.Fatalf("expect: nil, got: %v", err)
	}

	if err = validate(&Field{"time", "2016-3-04 01:01:01"}, IsTime()); err == nil {
		t.Fatalf("expect: %v, got: nil", err)
	}
	if err = validate(&Field{"time", "2016-3-4 01:01:01"}, IsTime()); err == nil {
		t.Fatalf("expect: %v, got: nil", err)
	}
	if err = validate(&Field{"time", "2016-03-4 01:01:01"}, IsTime()); err == nil {
		t.Fatalf("expect: %v, got: nil", err)
	}
	if err = validate(&Field{"time", "2016-03-04 24:01:01"}, IsTime()); err == nil {
		t.Fatalf("expect: %v, got: nil", err)
	}
	if err = validate(&Field{"time", "2016-03-04 23:59:60"}, IsTime()); err == nil {
		t.Fatalf("expect: %v, got: nil", err)
	}
}

func TestTimeAfter(t *testing.T) {
	var err error
	if err = validate(&Field{"time-after", time.Now().Add(time.Second)}, TimeAfter(time.Now())); err != nil {
		t.Fatalf("expect: nil, actual: %v", err)
	}
	if err = validate(&Field{"time-after", time.Now()}, TimeAfter("2018-09-09 00:00:00")); err != nil {
		t.Fatalf("expect: nil, actual: %v", err)
	}

	if err = validate(&Field{"time-after", time.Now()}, TimeAfter(time.Now().Add(time.Second))); err == nil {
		t.Fatalf("expect: %v, actual: nil", err)
	}
	if err = validate(&Field{"time", time.Now()}, TimeAfter("2090-02-03 00:00:00")); err == nil {
		t.Fatalf("expect: %v, actual: nil", err)
	}
}

func TestTimeBefore(t *testing.T) {
	var err error
	if err = validate(&Field{"time", time.Now()}, TimeBefore(time.Now().Add(time.Second))); err != nil {
		t.Fatalf("expect: nil, actual: %v", err)
	}
	if err = validate(&Field{"before-time", time.Now()}, TimeBefore("2050-02-03 00:00:00")); err != nil {
		t.Fatalf("expect: nil, actual: %v", err)
	}

	if err = validate(&Field{"before-time", time.Now().Add(time.Second)}, TimeBefore(time.Now())); err == nil {
		t.Fatalf("expect: %v, actual: nil", err)
	}
	if err = validate(&Field{"time", time.Now()}, TimeBefore("2018-02-03 00:00:00")); err == nil {
		t.Fatalf("expect: %v, actual: nil", err)
	}
}

func TestIsChinese(t *testing.T) {
	err := IsChinese()(&Field{Tag: "chinese", Value: reflect.ValueOf("中国").Interface()})
	if err != nil {
		t.Fatalf("expect nil, actual %v", err)
	}

	err = IsChinese()(&Field{Tag: "chinese", Value: reflect.ValueOf("chinese").Interface()})
	if err == nil {
		t.Fatalf("expect %v, actual nil", err)
	}

	err = IsChinese()(&Field{Tag: "chinese", Value: reflect.ValueOf("中国chinese").Interface()})
	if err == nil {
		t.Fatalf("expect %v, actual nil", err)
	}
}
