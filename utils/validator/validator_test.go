package validator

import (
	"testing"
)

func TestValidator(t *testing.T) {
	var val = NewValidator()
	ok := val.Validate(&Field{Tag: "age", Value: 10}, Min("100")).Ok()
	if ok {
		t.Errorf("expect: false, actual true")
	}
}

func TestParseValidatorAst(t *testing.T) {
	var expects = []string{
		"User.RawUint8",
		"User.RawUint8P",
		"User.RawInt8",
		"User.RawInt8P",
		"User.RawUint16",
		"User.RawUint16P",
		"User.RawInt16",
		"User.RawInt16P",
		"User.RawInt32",
		"User.RawInt32P",
		"User.RawUint32",
		"User.RawUint32P",
		"User.RawUint64",
		"User.RawUint64P",
		"User.RawInt64",
		"User.RawInt64P",
		"User.RawBool",
		"User.RawBoolP",
		"User.RawByte",
		"User.RawByteP",
		"User.RawBytes",
		"User.RawBytesP",
		"User.RawBytesPP",
		"User.RawBytesPPP",

		"User.VName",
		"User.VNameP",

		"User.VInfo.Address",
		"User.VInfo.Location.Latitude",
		"User.VInfo.Location.Longitude",
		"User.VInfoP.Address",
		"User.VInfoP.Location.Latitude",
		"User.VInfoP.Location.Longitude",

		"User.VSign",
		"User.VSignMulti",

		"User.VSlice",
		"User.MSign",
		"User.MSignP",
		"User.MRaw",
		"User.MRawP",

		"User.Location.Latitude",
		"User.Location.Longitude",
		"User.Info.Address",
		"User.Info.Location.Latitude",
		"User.Info.Location.Longitude",

		"User.VAction.VCall.CallNumber",
		"User.VAction.VCall.VVCall.VVCallNumber",
		"User.VActionP.VCallP.CallNumberP",
		"User.VActionP.VCallP.VVCallP.VVCallNumberP",
	}

	var linkeNames, err = ParseValidatorAst("User", "validator_file_test.go")
	if err != nil {
		t.Fatalf("expect: nil, actual:%v", err)
	}
	for i, linkName := range linkeNames {
		if expects[i] != linkName {
			t.Fatalf("expect: %v, actual:%v", expects[i], linkName)
		}
	}
	if len(linkeNames) != len(expects) {
		t.Fatalf("expect: %d, actual: %d", len(expects), len(linkeNames))
	}
}
