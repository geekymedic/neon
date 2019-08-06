package v9

import (
	"github.com/stretchr/testify/assert"
	"testing"
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

func TestCertNumber(t *testing.T) {
	var validator = &defaultValidator{}
	var arg = struct {
		CertNumber string `binding:"et_cert"`
	}{
		CertNumber: "92440606L402631312",
	}
	
	t.Run("", func(t *testing.T) {
		assert.Nil(t, validator.validate.Struct(arg))
	})
}