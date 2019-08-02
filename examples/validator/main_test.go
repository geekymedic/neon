package validator

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/geekymedic/neon/utils/validator"
)

//go:generate neon-cli validator generate --src-file=main.go --struct-name=Book --package-name=validator
type Book struct {
	Price     int    `json:"price"`
	Name      string `json:"name"`
	ExtraDesc Desc   `json:"extra_desc"`
	//Sign      map[string]Sign `json:"sign"`
}

type Desc struct {
	PublishingHouse string `json:"publishing_house"`
	Author          string `json:"author"`
}

type Sign struct {
	Dest       string
	AuthorName string
}

func TestMain1(t *testing.T) {
	http.HandleFunc("/brook/add", func(writer http.ResponseWriter, request *http.Request) {
		var arg Book
		err := json.NewDecoder(request.Body).Decode(&arg)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte(err.Error()))
			return
		}
		valBook := NewBookValidator(&arg)
		val := validator.NewValidator()
		if err := val.Validate(valBook.Price, validator.Min("2000")).Err(); err != nil {
			writer.Write([]byte(err.Error()))
			return
		}
		json.NewEncoder(writer).Encode(arg)
	})
	go func() {
		err := http.ListenAndServe(":1349", nil)
		if err != nil {
			t.Errorf("expect: nil, actual: %v", err)
		}
	}()
	var arg = Book{
		Price: 10200,
	}
	buf, _ := json.Marshal(arg)
	resp, err := http.Post("http://localhost:1349/brook/add", "", bytes.NewReader(buf))
	if err != nil {
		t.Errorf("expect: nil, actual: %v", err)
	}
	buf, _ = ioutil.ReadAll(resp.Body)
	t.Log(string(buf))
}
