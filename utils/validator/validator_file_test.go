package validator

const (
	Zero = iota
	FirstNumber
	SecondNumber
)

//go:generate neon-cli validator generate --src-file=validator_file_test.go --struct-name=Book --package-name=controll
type Book struct {
	Name  string `json:"string"`
	Price float64
	Fact  Factory `json:"factory"`
}

type Factory struct {
	Address string `json:"address"`
}

//go:generate neon-cli validator generate --src-file=validator_file_test.go --struct-name=User --package-name=validator --save-file=/tmp/validator_ext_test.go
type User struct {
	RawUint8    uint8
	RawUint8P   *uint8
	RawInt8     int8
	RawInt8P    *int8
	RawUint16   uint16
	RawUint16P  *uint16
	RawInt16    int16
	RawInt16P   *int16
	RawInt32    int32
	RawInt32P   *int32
	RawUint32   uint32
	RawUint32P  *uint32
	RawUint64   uint64
	RawUint64P  *uint64
	RawInt64    int64
	RawInt64P   *int64
	RawBool     bool
	RawBoolP    *bool
	RawByte     byte
	RawByteP    *byte
	RawBytes    []byte
	RawBytesP   []*byte
	RawBytesPP  *[]byte
	RawBytesPPP *[]*byte

	VName      string // ast.Ident
	VNameP     *string
	VInfo      Info `json:"tag_info"` // ast.TypeSpec
	VInfoP     *Info
	VSign      []Sign `json:"tag_sign"`
	VSignMulti [][][]Sign
	VSlice     []int
	MSign      map[string]Sign
	MSignP     map[string]*Sign
	MRaw       map[string]int
	MRawP      map[string]*int
	Location   `json:"tag_location"`
	*Info
	VAction struct {
		// ast.Ident
		VCall struct {
			CallNumber int `json:"tag_call_number"`
			VVCall     struct {
				VVCallNumber int `json:"tag_vvcall_number"`
			} `json:"tag_vvcall"`
		} `json:"tag_vcall"`
	} `json:"tag_vaction"`
	VActionP *struct {
		// ast.Ident
		VCallP *struct {
			CallNumberP *int
			VVCallP     *struct {
				VVCallNumberP *int
			}
		}
	}
}

type Info struct {
	Address  string
	Location Location
}

type Location struct {
	Latitude  string
	Longitude string
}

type Sign struct {
	Dest   string
	Author string
}

func newChannel() <-chan int {
	ch := make(chan int, 10)
	go func() {
		<-ch
	}()
	return ch
}
