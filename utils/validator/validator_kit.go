package validator

import (
	"context"
	"fmt"
	"github.com/geekymedic/neon/errors"
	valkitv8 "github.com/geekymedic/neon/utils/validator/v8"
	"reflect"
	"time"
)

var validatorVersion = v8
var ErrUnImplement = errors.NewStackError("unimplement")

const (
	v8 = "v8"
	v9 = "v9"
)

// TODO add v9 support
func SetValidatorVersion(version string) {
	if version != v8 {
		panic("version must be v8 or v9")
	}
	validatorVersion = version
}

func IsEq(param string) OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, param, valkitv8.IsEq, "IsEq")
		}
		return ErrUnImplement
	}
}

func Len(param string) OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, param, valkitv8.HasLengthOf, "Len")
		}
		return ErrUnImplement
	}
}

// min(param) <= field.Value
func Min(param string) OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, param, valkitv8.HasMinOf, "Min")
		}
		return ErrUnImplement
	}
}

// field.Value <= max(param)
func Max(param string) OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, param, valkitv8.HasMaxOf, "Max")
		}
		return ErrUnImplement
	}
}

func Eq(param string) OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, param, valkitv8.IsEq, "Eq")
		}
		return ErrUnImplement
	}
}

func Ne(param string) OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, param, valkitv8.IsNe, "IsNe")
		}
		return ErrUnImplement
	}
}

func Lt(param string) OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, param, valkitv8.IsLt, "Lt")
		}
		return ErrUnImplement
	}
}

func Lte(param string) OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, param, valkitv8.IsLte, "IsLet")
		}
		return ErrUnImplement
	}
}

func Gt(param string) OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, param, valkitv8.IsGt, "IsGt")
		}
		return ErrUnImplement
	}
}

func Gte(param string) OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, param, valkitv8.IsGte, "IsGte")
		}
		return ErrUnImplement
	}
}

func IsAlpha() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsAlpha, "IsAlpha")
		}
		return ErrUnImplement
	}
}

func IsAlphanum() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsAlphanum, "IsAlphanum")
		}
		return ErrUnImplement
	}
}

func IsNumeric() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsNumeric, "IsNumeric")
		}
		return ErrUnImplement
	}
}

func IsNumber() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsNumber, "IsNumber")
		}
		return ErrUnImplement
	}
}

func IsChinese() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsChinese, "IsChinese")
		}
		return ErrUnImplement
	}
}

func IsNickName() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsNickName, "IsNickName")
		}
		return ErrUnImplement
	}
}

func Hexadecimal() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsHexadecimal, "IsHexadecimal")
		}
		return ErrUnImplement
	}
}

func Hexcolor() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsHEXColor, "IsHEXColor")
		}
		return ErrUnImplement
	}
}

func IsRGB() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsRGB, "IsRGB")
		}
		return ErrUnImplement
	}
}

func IRGBA() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsRGBA, "IsRGBA")
		}
		return ErrUnImplement
	}
}

func IsHSL() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsHSL, "IsHSL")
		}
		return ErrUnImplement
	}
}

func IsHSLA() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsHSLA, "IsHSLA")
		}
		return ErrUnImplement
	}
}

func IsEmail() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsEmail, "IsEmail")
		}
		return ErrUnImplement
	}
}

func IsURL() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsURL, "IsURL")
		}
		return ErrUnImplement
	}
}

func IsURI() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsURI, "IsURI")
		}
		return ErrUnImplement
	}
}

func IsBase64() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsBase64, "IsBase64")
		}
		return ErrUnImplement
	}
}

func Contains(param string) OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, param, valkitv8.Contains, "Contains")
		}
		return ErrUnImplement
	}
}

func ContainsAny(param string) OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, param, valkitv8.ContainsAny, "ContainsAny")
		}
		return ErrUnImplement
	}
}

func ContainsRune(param string) OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, param, valkitv8.ContainsRune, "ContainsRune")
		}
		return ErrUnImplement
	}
}

func Excludes(param string) OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, param, valkitv8.Excludes, "Excludes")
		}
		return ErrUnImplement
	}
}

func ExcludesAll(param string) OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, param, valkitv8.ExcludesAll, "ExcludesAll")
		}
		return ErrUnImplement
	}
}

func ExcludesRune(param string) OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, param, valkitv8.ExcludesRune, "ExcludesRune")
		}
		return ErrUnImplement
	}
}

func IsISBN() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsISBN, "IsISBN")
		}
		return ErrUnImplement
	}
}

func IsISBN10() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsISBN10, "IsISBN10")
		}
		return ErrUnImplement
	}
}

func IsISBN13() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsISBN13, "IsISBN13")
		}
		return ErrUnImplement
	}
}

func IsUUID() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsUUID, "IsUUID")
		}
		return ErrUnImplement
	}
}

func IsUUID3() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsUUID3, "IsUUID3")
		}
		return ErrUnImplement
	}
}

func IsUUID4() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsUUID4, "IsUUID4")
		}
		return ErrUnImplement
	}
}

func IsUUID5() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsUUID5, "IsUUID5")
		}
		return ErrUnImplement
	}
}

func IsASCII() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsASCII, "IsASCII")
		}
		return ErrUnImplement
	}
}

func IsPrintableASCII() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsPrintableASCII, "IsPrintableASCII")
		}
		return ErrUnImplement
	}
}

func IsMultiByteCharacter() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.HasMultiByteCharacter, "IsMultiByteCharacter")
		}
		return ErrUnImplement
	}
}

func IsDataURI() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsDataURI, "IsDataURI")
		}
		return ErrUnImplement
	}
}

func IsLatitude() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsLatitude, "IsLatitude")
		}
		return ErrUnImplement
	}
}

func IsLongitude() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsLongitude, "IsLongitude")
		}
		return ErrUnImplement
	}
}

func IsSSN() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsSSN, "IsSSN")
		}
		return ErrUnImplement
	}
}

func IsIPv4() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsIPv4, "IsIPv4")
		}
		return ErrUnImplement
	}
}

func IsIPv6() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsIPv6, "IsIPv6")
		}
		return ErrUnImplement
	}
}

func IsIP() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsIP, "IsIP")
		}
		return ErrUnImplement
	}
}

func IsCIDRv4() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsCIDRv4, "IsCIDRv4")
		}
		return ErrUnImplement
	}
}

func IsCIDRv6() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsCIDRv6, "IsCIDRv6")
		}
		return ErrUnImplement
	}
}

func IsCIDR() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsCIDR, "IsCIDR")
		}
		return ErrUnImplement
	}
}

func IsTCP4AddrResolvable() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsTCP4AddrResolvable, "IsTCP4AddrResolvable")
		}
		return ErrUnImplement
	}
}

func IsTCP6AddrResolvable() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsTCP6AddrResolvable, "IsTCP6AddrResolvable")
		}
		return ErrUnImplement
	}
}

func IsTCPAddrResolvable() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsTCPAddrResolvable, "IsTCPAddrResolvable")
		}
		return ErrUnImplement
	}
}

func IsUDP4AddrResolvable() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsUDP4AddrResolvable, "IsUDP4AddrResolvable")
		}
		return ErrUnImplement
	}
}

func IsUDP6AddrResolvable() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsUDP6AddrResolvable, "IsUDP6AddrResolvable")
		}
		return ErrUnImplement
	}
}

func IsUDPAddrResolvable() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsUDPAddrResolvable, "IsUDPAddrResolvable")
		}
		return ErrUnImplement
	}
}

func IsIP4AddrResolvable() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsIP4AddrResolvable, "IsIP4AddrResolvable")
		}
		return ErrUnImplement
	}
}

func IsIP6AddrResolvable() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsIP6AddrResolvable, "IsIP6AddrResolvable")
		}
		return ErrUnImplement
	}
}

func IsIPAddrResolvable() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsIPAddrResolvable, "IsIPAddrResolvable")
		}
		return ErrUnImplement
	}
}

func IsUnixAddrResolvable() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsUnixAddrResolvable, "IsUnixAddrResolvable")
		}
		return ErrUnImplement
	}
}

func IsMAC() OpFn {
	return func(field *Field) error {
		if validatorVersion == v8 {
			return wrapV8Fn(context.Background(), field, "", valkitv8.IsMAC, "IsMAC")
		}
		return ErrUnImplement
	}
}

// UTC time
func IsTime() OpFn {
	return func(field *Field) error {
		return wrapV8Fn(context.Background(), field, "", isTime, "IsTime")
	}
}

// UTC time
// TimeAfter reports whether the field.time instant t is after arg.time.
func TimeAfter(arg interface{}) OpFn {
	return func(field *Field) error {
		var argTime, fieldTime time.Time
		var err error
		switch typ := arg.(type) {
		case time.Time:
			argTime = typ
		case *time.Time:
			argTime = *typ
		case string:
			argTime, err = time.ParseInLocation("2006-01-02 15:04:05", typ, time.Local)
		case *string:
			argTime, err = time.ParseInLocation("2006-01-02 15:04:05", *typ, time.Local)
		default:
			err = errors.NewStackError(fmt.Sprintf("tag: %v, op: %v, real type: %v, real value: %v, param: %v", field.Tag, "TimeBefore",
				reflect.TypeOf(field.Value).String(), field.Value, arg))
			return err
		}

		switch typ := field.Value.(type) {
		case time.Time:
			fieldTime = typ
		case *time.Time:
			fieldTime = *typ
		default:
			err = errors.NewStackError(fmt.Sprintf("tag: %v, op: %v, real type: %v, real value: %v, param: %v", field.Tag, "TimeBefore",
				reflect.TypeOf(field.Value).String(), field.Value, arg))
			return err
		}

		if err != nil {
			return err
		}
		if fieldTime.After(argTime) {
			return nil
		}
		err = errors.NewStackError(fmt.Sprintf("tag: %v, op: %v, real type: %v, real value: %v, param: %v", field.Tag, "TimeBefore",
			reflect.TypeOf(field.Value).String(), field.Value, arg))
		return err
	}
}

// UTC time
func TimeBefore(arg interface{}) OpFn {
	return func(field *Field) error {
		var argTime, fieldTime time.Time
		var err error
		switch typ := arg.(type) {
		case time.Time:
			argTime = typ
		case *time.Time:
			argTime = *typ
		case string:
			argTime, err = time.ParseInLocation("2006-01-02 15:04:05", typ, time.Local)
		case *string:
			argTime, err = time.ParseInLocation("2006-01-02 15:04:05", *typ, time.Local)
		default:
			err = errors.NewStackError(fmt.Sprintf("tag: %v, op: %v, real type: %v, real value: %v, param: %v", field.Tag, "TimeBefore",
				reflect.TypeOf(field.Value).String(), field.Value, arg))
			return err
		}

		switch typ := field.Value.(type) {
		case time.Time:
			fieldTime = typ
		case *time.Time:
			fieldTime = *typ
		default:
			err = errors.NewStackError(fmt.Sprintf("tag: %v, op: %v, real type: %v, real value: %v, param: %v", field.Tag, "TimeBefore",
				reflect.TypeOf(field.Value).String(), field.Value, arg))
			return err
		}
		if fieldTime.Before(argTime) {
			return nil
		}
		err = errors.NewStackError(fmt.Sprintf("tag: %v, op: %v, real type: %v, real value: %v, param: %v", field.Tag, "TimeBefore",
			reflect.TypeOf(field.Value).String(), field.Value, arg))
		return err
	}
}

func isTime(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	_, err := time.ParseInLocation("2006-01-02 15:04:05", value.String(), time.Local)
	return err == nil
}

func RegisterValidation(key string, fn valkitv8.Func) error {
	return valkitv8.RegisterValidation(key, fn)
}

func wrapV8Fn(ctx context.Context, field *Field, param string, fn func(context.Context, interface{}, interface{}) bool, fnName string) error {
	ok := fn(ctx, field.Value, param)
	if ok {
		return nil
	}
	typ := reflect.TypeOf(field.Value)
	err := errors.NewStackError(fmt.Sprintf("tag: %v, op: %v, real type: %v, real value: %v, param: %v", field.Tag, fnName, typ.String(), field.Value, param))
	return errors.By(err)
}
