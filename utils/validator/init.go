package validator

import (
	"reflect"

	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v8"
)

func init() {
	binding.Validator.Engine().(*validator.Validate).RegisterValidation("chinese", func(v *validator.Validate, topStruct reflect.Value, currentStruct reflect.Value, field reflect.Value, fieldtype reflect.Type, fieldKind reflect.Kind, param string) bool {
		return IsChinese()(&Field{Tag: "chinese", Value: field.Interface()}) == nil
	})

	binding.Validator.Engine().(*validator.Validate).RegisterValidation("nickname", func(v *validator.Validate, topStruct reflect.Value, currentStruct reflect.Value, field reflect.Value, fieldtype reflect.Type, fieldKind reflect.Kind, param string) bool {
		return IsNickName()(&Field{Tag: "nickname", Value: field.Interface()}) == nil
	})
}
