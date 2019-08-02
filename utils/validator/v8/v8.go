package v8

import (
	"context"
	"fmt"
	"reflect"
	"unicode"
	"unicode/utf8"

	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v8"
)

var (
	ignoreValue reflect.Value
	ignoreKind  reflect.Kind
	ignoreType  reflect.Type
	ignoreParam string
)

type Func = validator.Func
type Validate = validator.Validate

// IsEq is the validation function for validating if the current field's value is equal to the param's value.
func IsEq(_ context.Context, field interface{}, param interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsEq(nil, ignoreValue, ignoreValue, value, nil, value.Kind(), convertParam(param))
	return ok
}

// IsGt is the validation function for validating if the current field's value is greater than the param's value.
func IsNe(ctx context.Context, field interface{}, param interface{}) bool {
	return !IsEq(ctx, field, param)
}

// HasLengthOf is the validation function for validating if the current field's value is equal to the param's value.
func HasLengthOf(_ context.Context, field interface{}, param interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.HasLengthOf(nil, ignoreValue, ignoreValue, value, nil, value.Kind(), convertParam(param))
	return ok
}

// HasMinOf is the validation function for validating if the current field's value is greater than or equal to the param's value.
func HasMinOf(_ context.Context, field interface{}, param interface{}) bool {
	value := reflect.ValueOf(field)
	typ := reflect.TypeOf(field)
	ok := validator.HasMinOf(nil, ignoreValue, ignoreValue, value, typ, value.Kind(), convertParam(param))
	return ok
}

// HasMaxOf is the validation function for validating if the current field's value is less than or equal to the param's value.
func HasMaxOf(_ context.Context, field interface{}, param interface{}) bool {
	value := reflect.ValueOf(field)
	typ := reflect.TypeOf(field)
	ok := validator.HasMaxOf(nil, ignoreValue, ignoreValue, value, typ, value.Kind(), convertParam(param))
	return ok
}

// IsLt is the validation function for validating if the current field's value is less than the param's value.
func IsLt(_ context.Context, field interface{}, param interface{}) bool {
	value := reflect.ValueOf(field)
	typ := reflect.TypeOf(field)
	ok := validator.IsLt(nil, ignoreValue, ignoreValue, value, typ, value.Kind(), convertParam(param))
	return ok
}

// IsLte is the validation function for validating if the current field's value is less than or equal to the param's value.
func IsLte(_ context.Context, field interface{}, param interface{}) bool {
	value := reflect.ValueOf(field)
	typ := reflect.TypeOf(field)
	ok := validator.IsLte(nil, ignoreValue, ignoreValue, value, typ, value.Kind(), convertParam(param))
	return ok
}

// IsGt is the validation function for validating if the current field's value is greater than the param's value.
func IsGt(_ context.Context, field interface{}, param interface{}) bool {
	value := reflect.ValueOf(field)
	typ := reflect.TypeOf(field)
	ok := validator.IsGt(nil, ignoreValue, ignoreValue, value, typ, value.Kind(), convertParam(param))
	return ok
}

// IsGte is the validation function for validating if the current field's value is greater than or equal to the param's value.
func IsGte(_ context.Context, field interface{}, param interface{}) bool {
	value := reflect.ValueOf(field)
	typ := reflect.TypeOf(field)
	ok := validator.IsGte(nil, ignoreValue, ignoreValue, value, typ, value.Kind(), convertParam(param))
	return ok
}

// IsAlpha is the validation function for validating if the current field's value is a valid alpha value.
func IsAlpha(_ context.Context, field interface{}, param interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsAlpha(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsAlphanum is the validation function for validating if the current field's value is a valid alphanumeric value.
func IsAlphanum(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsAlphanum(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsNumeric is the validation function for validating if the current field's value is a valid numeric value.
func IsNumeric(_ context.Context, field interface{}, param interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsNumeric(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsNumber is the validation function for validating if the current field's value is a valid number.
func IsNumber(_ context.Context, field interface{}, param interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsNumber(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

func IsChinese(_ context.Context, field interface{}, param interface{}) bool {
	value := reflect.ValueOf(field)
	if value.Kind() != reflect.String {
		return false
	}
	for _, runeValue := range value.String() {
		if !unicode.Is(unicode.Han, runeValue) {
			return false
		}
	}
	return true
}

func IsNickName(_ context.Context, field interface{}, param interface{}) bool {
	value := reflect.ValueOf(field)
	if value.Kind() != reflect.String {
		return false
	}
	strValue := value.String()
	for i, w := 0, 0; i < len(strValue); i += w {
		runeValue, width := utf8.DecodeRuneInString(strValue[i:])
		w = width
		if unicode.Is(unicode.Han, runeValue) {
			//fmt.Printf("%#U 汉字\n", runeValue)
			continue
		}
		if runeValue == rune('_') {
			//fmt.Printf("%#U 下划线\n", runeValue)
			continue
		}

		if (runeValue >= rune('A') && runeValue <= rune('Z')) || (runeValue >= rune('a') && runeValue <= rune('z')) {
			//fmt.Printf("%#U 英文\n", runeValue)
			continue
		}

		if runeValue >= rune('0') && runeValue <= rune('9') {
			//fmt.Printf("%#U 数字\n", runeValue)
			continue
		}

		return false
	}
	return true
}

// IsNumber is the validation function for validating if the current field's value is a valid number.
func IsHexadecimal(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsHexadecimal(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsHEXColor is the validation function for validating if the current field's value is a valid HEX color.
func IsHEXColor(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsHEXColor(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsRGB is the validation function for validating if the current field's value is a valid RGB color.
func IsRGB(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsRGB(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsRGBA is the validation function for validating if the current field's value is a valid RGBA color.
func IsRGBA(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsRGBA(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsHSL is the validation function for validating if the current field's value is a valid HSL color.
func IsHSL(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsHSL(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsHSLA is the validation function for validating if the current field's value is a valid HSLA color.
func IsHSLA(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsHSLA(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsEmail is the validation function for validating if the current field's value is a valid email address.
func IsEmail(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsEmail(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsURL is the validation function for validating if the current field's value is a valid URL.
func IsURL(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsURL(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsURI is the validation function for validating if the current field's value is a valid URI.
func IsURI(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsURI(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsBase64 is the validation function for validating if the current field's value is a valid base 64.
func IsBase64(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsBase64(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// Contains is the validation function for validating that the field's value contains the text specified within the param.
func Contains(_ context.Context, field interface{}, param interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.Contains(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, convertParam(param))
	return ok
}

// ContainsAny is the validation function for validating that the field's value contains any of the characters specified within the param.
func ContainsAny(_ context.Context, field interface{}, param interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.ContainsAny(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, convertParam(param))
	return ok
}

// ContainsRune is the validation function for validating that the field's value contains the rune specified within the param.
func ContainsRune(_ context.Context, field interface{}, param interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.ContainsRune(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, convertParam(param))
	return ok
}

// Excludes is the validation function for validating that the field's value does not contain the text specified within the param.
func Excludes(_ context.Context, field interface{}, param interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.Excludes(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, convertParam(param))
	return ok
}

// ExcludesAll is the validation function for validating that the field's value does not contain any of the characters specified within the param.
func ExcludesAll(_ context.Context, field interface{}, param interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.ExcludesAll(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, convertParam(param))
	return ok
}

// ExcludesRune is the validation function for validating that the field's value does not contain the rune specified within the param.
func ExcludesRune(_ context.Context, field interface{}, param interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.ExcludesRune(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, convertParam(param))
	return ok
}

// IsISBN is the validation function for validating if the field's value is a valid v10 or v13 ISBN.
func IsISBN(ctx context.Context, field interface{}, param interface{}) bool {
	ok := IsISBN10(ctx, field, param) || IsISBN13(ctx, field, param)
	return ok
}

// IsISBN10 is the validation function for validating if the field's value is a valid v10 ISBN.
func IsISBN10(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsISBN10(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsISBN13 is the validation function for validating if the field's value is a valid v13 ISBN.
func IsISBN13(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsISBN13(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsUUID is the validation function for validating if the field's value is a valid UUID of any version.
func IsUUID(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsUUID(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsUUID3 is the validation function for validating if the field's value is a valid v3 UUID.
func IsUUID3(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsUUID3(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsUUID4 is the validation function for validating if the field's value is a valid v4 UUID.
func IsUUID4(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsUUID4(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsUUID5 is the validation function for validating if the field's value is a valid v5 UUID.
func IsUUID5(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsUUID5(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsASCII is the validation function for validating if the field's value is a valid ASCII character.
func IsASCII(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsASCII(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsPrintableASCII is the validation function for validating if the field's value is a valid printable ASCII character.
func IsPrintableASCII(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsPrintableASCII(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// HasMultiByteCharacter is the validation function for validating if the field's value has a multi byte character.
func HasMultiByteCharacter(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.HasMultiByteCharacter(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsDataURI is the validation function for validating if the field's value is a valid data URI.
func IsDataURI(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsDataURI(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsLatitude is the validation function for validating if the field's value is a valid latitude coordinate.
func IsLatitude(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsLatitude(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsLongitude is the validation function for validating if the field's value is a valid longitude coordinate.
func IsLongitude(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsLongitude(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsSSN is the validation function for validating if the field's value is a valid SSN.
func IsSSN(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsSSN(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsIPv4 is the validation function for validating if a value is a valid v4 IP address.
func IsIPv4(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsIPv4(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsIPv6 is the validation function for validating if the field's value is a valid v6 IP address.
func IsIPv6(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsIPv6(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsIP is the validation function for validating if the field's value is a valid v4 or v6 IP address.
func IsIP(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsIP(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsCIDRv4 is the validation function for validating if the field's value is a valid v4 CIDR address.
func IsCIDRv4(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsCIDRv4(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsCIDRv6 is the validation function for validating if the field's value is a valid v6 CIDR address.
func IsCIDRv6(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsCIDRv6(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsCIDR is the validation function for validating if the field's value is a valid v4 or v6 CIDR address.
func IsCIDR(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsCIDR(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsTCP4AddrResolvable is the validation function for validating if the field's value is a resolvable tcp4 address.
func IsTCP4AddrResolvable(ctx context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsTCP4AddrResolvable(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsTCP6AddrResolvable is the validation function for validating if the field's value is a resolvable tcp6 address.
func IsTCP6AddrResolvable(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsTCP6AddrResolvable(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsTCPAddrResolvable is the validation function for validating if the field's value is a resolvable tcp address.
func IsTCPAddrResolvable(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsTCPAddrResolvable(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsUDP4AddrResolvable is the validation function for validating if the field's value is a resolvable udp4 address.
func IsUDP4AddrResolvable(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsUDP4AddrResolvable(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsUDP6AddrResolvable is the validation function for validating if the field's value is a resolvable udp6 address.
func IsUDP6AddrResolvable(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsUDP6AddrResolvable(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsUDPAddrResolvable is the validation function for validating if the field's value is a resolvable udp address.
func IsUDPAddrResolvable(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsUDPAddrResolvable(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsIP4AddrResolvable is the validation function for validating if the field's value is a resolvable ip4 address.
func IsIP4AddrResolvable(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsIP4AddrResolvable(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsIP6AddrResolvable is the validation function for validating if the field's value is a resolvable ip6 address.
func IsIP6AddrResolvable(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsIP6AddrResolvable(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsIPAddrResolvable is the validation function for validating if the field's value is a resolvable ip address.
func IsIPAddrResolvable(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsIPAddrResolvable(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsUnixAddrResolvable is the validation function for validating if the field's value is a resolvable unix address.
func IsUnixAddrResolvable(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsUnixAddrResolvable(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

// IsMAC is the validation function for validating if the field's value is a valid MAC address.
func IsMAC(_ context.Context, field interface{}, _ interface{}) bool {
	value := reflect.ValueOf(field)
	ok := validator.IsMAC(nil, ignoreValue, ignoreValue, value, ignoreType, ignoreKind, ignoreParam)
	return ok
}

func RegisterValidation(key string, fn validator.Func) error {
	return binding.Validator.Engine().(*validator.Validate).RegisterValidation(key, fn)
}

func convertParam(param interface{}) string {
	if str, ok := param.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", param)
}
