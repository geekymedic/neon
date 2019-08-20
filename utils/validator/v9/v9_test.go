package v9

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestIdentify(t *testing.T) {
	var validator = &defaultValidator{}
	var arg = struct {
		IdentfyId string `binding:"et_identity"`
	}{
		IdentfyId: "440881188107238765",
	}
	validator.lazyinit()
	assert.Nil(t, validator.validate.Struct(arg))

	arg.IdentfyId = "x440881188107238765"
	assert.NotNil(t, validator.validate.Struct(arg))
}

func TestPhoneNumber(t *testing.T) {
	var validator = &defaultValidator{}
	var arg = struct {
		Phone string `binding:"et_phone"`
	}{
		Phone: "13760049089",
	}
	validator.lazyinit()
	assert.Nil(t, validator.validate.Struct(arg))

	arg.Phone = "10798430128"
	assert.NotNil(t, validator.validate.Struct(arg))
}

func TestExContains(t *testing.T) {
	var validator = &defaultValidator{}
	validator.lazyinit()
	t.Run("int", func(t *testing.T) {
		var arg = struct {
			Age int `binding:"et_contains=1-2"`
		}{
			Age: 1,
		}
		assert.Nil(t, validator.validate.Struct(arg))
		arg.Age = 10
		assert.NotNil(t, validator.validate.Struct(arg))
	})

	t.Run("string", func(t *testing.T) {
		var arg = struct {
			Name string `binding:"et_contains=abc-cda"`
		}{
			Name: "abc",
		}
		assert.Nil(t, validator.validate.Struct(arg))
		arg.Name = "123"
		assert.NotNil(t, validator.validate.Struct(arg))
	})

	t.Run("bool", func(t *testing.T) {

	})
}

//func TestCertNumber(t *testing.T) {
//	var validator = &defaultValidator{}
//	var arg = struct {
//		CertNumber string `binding:"et_cert"`
//	}{
//		CertNumber: "92440606L402631312",
//	}
//
//	t.Run("", func(t *testing.T) {
//		assert.Nil(t, validator.validate.Struct(arg))
//	})
//}

func TestJson(t *testing.T) {
	var validator = &defaultValidator{}

	validator.lazyinit()

	t.Run("ok", func(t *testing.T) {
		var arg = struct {
			CertNumber string `binding:"et_json"`
		}{
			CertNumber: "{}",
		}
		assert.Nil(t, validator.validate.Struct(arg))
	})

	t.Run("fail", func(t *testing.T) {
		var arg = struct {
			CertNumber string `binding:"et_json"`
		}{
			CertNumber: "{''}",
		}
		assert.NotNil(t, validator.validate.Struct(arg))
	})
}

type ArgTime time.Time

func (t *ArgTime) UnmarshalJSON(buf []byte) error {
	var tm, err = time.Parse("2006-01-02 15:04:05", "2019-03-29 15:04:05")
	if err != nil {
		return err
	}
	*t = ArgTime(tm)
	return nil
}

func (t *ArgTime) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(*t).Format("2006-01-02 15:04:05")), nil
}

func TestTime(t *testing.T) {
	var validator = &defaultValidator{}
	validator.lazyinit()

	type Arg struct {
		StartTime ArgTime
	}

	t.Run("", func(t *testing.T) {
		var arg = Arg{

		}
		//validator.validate.Struct(arg)
		err := json.Unmarshal([]byte(`{"StartTime":"2019-03-29 15:04:05"}`), &arg)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(arg)

		var p *int
		buf, err := json.Marshal(p)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("-->", string(buf))

		var pp *Arg
		buf, err = json.Marshal(pp135004854809348715
		)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("-->", string(buf))
	})
}
