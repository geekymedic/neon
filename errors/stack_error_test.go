package errors

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestErrChain_Error(t *testing.T) {
	t.Run("", func(t *testing.T) {
		var errc error = &errChain{err: errors.New("hello word"), msg: "chain1"}
		t.Log(errc)
	})
}

func TestWithMessageEx(t *testing.T) {
	t.Parallel()

	t.Run("simple test", func(t *testing.T) {
		expect := "fail to open the file"
		err := WithMessage(errors.New(expect), "%s", "helloword")
		if err == nil {
			t.Errorf("expect: %v, got: %v", err, nil)
		}
		if strings.Contains(err.Error(), expect) == false {
			t.Errorf("expect: %v, got: not found", expect)
		}
	})

	t.Run("create a new stack error with a std error", func(t *testing.T) {
		expect := "fail to open the file"
		err := WithMessage(errors.New("fail to open the file"), "%s %d", "helloword", 1000)
		if err == nil {
			t.Errorf("expect: %v, got: %v", err, nil)
		}
		ret := strings.Contains(err.Error(), expect)
		if ret == false {
			t.Errorf("expect: %v, got: %v", true, false)
		}
	})

	t.Run("create a new stack error with a stack error", func(t *testing.T) {
		expect := "fail to open the file"
		err := WithMessage(NewStackError("fail to open the file"), "%s %d", "helloword", 1000)
		if err == nil {
			t.Errorf("expect: %v, got: %v", err, nil)
		}
		ret := strings.Contains(err.Error(), expect)
		if ret == false {
			t.Errorf("expect: %v, got: %v", true, false)
		}
	})

	t.Run("create a new stack error with multi levels stack error", func(t *testing.T) {
		expect := "fail to open the file"
		var chains []string
		for i := 0; i < 10; i++ {
			chains = append(chains, fmt.Sprintf("chain%v", i))
		}
		err := WithMessage(errors.New(expect), "%s", "last")
		for i := 0; i < 10; i++ {
			err = WithMessage(err, "%s", chains[i])
		}

		if strings.Contains(err.Error(), expect) == false {
			t.Errorf("expect: %v, got: not found", expect)
		}
		for i := 0; i < 10; i++ {
			if strings.Contains(err.Error(), chains[0]) == false {
				t.Errorf("expect: %v, got: not found", chains[i])
			}
		}
	})
}

func TestFormatEx(t *testing.T) {
	t.Run("create a stack error be format", func(t *testing.T) {
		err := Format("%s", "timeout")
		if strings.Contains(err.Error(), "timeout") == false {
			t.Errorf("expect: %v, got: not found", "timeout")
		}
	})
}

func TestByEx(t *testing.T) {
	t.Run("create a nil error", func(t *testing.T) {
		err := By(nil)
		if err != nil {
			t.Errorf("expect: %v, got: %v", nil, err)
		}
	})

	t.Run("create a stack error with std error", func(t *testing.T) {
		argErr := errors.New("fail to open the file")
		err := By(errors.New("fail to open the file"))
		if err == nil {
			t.Errorf("expect: %v, got: %v", argErr, err)
		}
		t.Log(err)
	})

	t.Run("create a stack error with stack error", func(t *testing.T) {
		argErr := NewStackError("fail to open file")
		err := By(argErr)
		if err == nil {
			t.Errorf("expect: %v, got: %v", argErr, err)
		}
		t.Log(err.Error())
	})
}

func TestWrap(t *testing.T) {
	t.Run("err type not equal StackError", func(t *testing.T) {
		expectOutput := []string{
			`fail to open the file`,
			`github.com/geekymedic/neon/errors.Wrap`,
		}
		err := errors.New("fail to open the file")
		err = Wrap(err)
		for _, expect := range expectOutput {
			if !strings.Contains(err.Error(), expect) {
				t.Fatalf("expect: %v, got: %v", expect, err.Error())
			}
		}
	})

	t.Run("err type equal StackError", func(t *testing.T) {
		expectOutput := []string{
			`github.com/geekymedic/neon/errors.TestWrap.func2`,
		}
		err := NewStackError("fail to open file")
		err = Wrap(err)
		for _, expect := range expectOutput {
			if !strings.Contains(err.Error(), expect) {
				t.Fatalf("expect: %v, got: %v", expect, err.Error())
			}
		}
	})

	t.Run("err type equal StackError and mult call", func(t *testing.T) {
		expectOutput := []string{
			`fail to call a function`,
			`github.com/geekymedic/neon/errors.stack_call_a`,
			`github.com/geekymedic/neon/errors.stack_call_b`,
			`github.com/geekymedic/neon/errors.stack_call_c`,
		}
		err := stack_call_c()
		for _, expect := range expectOutput {
			if !strings.Contains(err.Error(), expect) {
				t.Fatalf("expect: %v, got: %v", expect, err.Error())
			}
		}
	})

	t.Run("err type equal StackError and mult call mannully", func(t *testing.T) {
		expectOutput := []string{
			`fail to open the file`,
			`github.com/geekymedic/neon/errors.TestWrap.func4`,
			`github.com/geekymedic/neon/errors/stack_error_test.go`,
		}
		err := NewStackError("fail to open the file")
		for _, expect := range expectOutput {
			if !strings.Contains(err.Error(), expect) {
				t.Fatalf("expect: %v, got: %v", expect, err.Error())
			}
		}
		err = Wrap(err)
		err = Wrap(err)
		for _, expect := range expectOutput {
			if !strings.Contains(err.Error(), expect) {
				t.Fatalf("expect: %v, got: %v", expect, err.Error())
			}
		}
	})

	t.Run("shortMsg", func(t *testing.T) {
		err := Wrap(errors.New("stack One"))
		t.Log(ShortMsg(err))
	})
}

//go:noinline
func stack_call_a() error {
	err := NewStackError("fail to call a function")
	return err
}

//go:noinline
func stack_call_b() error {
	err := stack_call_a()
	return Wrap(err)
}

//go:noinline
func stack_call_c() error {
	return stack_call_b()
}
