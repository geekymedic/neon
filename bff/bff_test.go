package bff

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/geekymedic/neon"
	"github.com/geekymedic/neon/logger"

	"github.com/stretchr/testify/assert"

	_ "github.com/geekymedic/neon/plugin/metrics"
)

func TestEngine(t *testing.T) {
	l, err := net.Listen("tcp", ":8080")
	assert.Nil(t, err)
	go func() {
		_engine.POST("/test", HttpHandler(func(state *State) {
			var id = struct {
				Id string `json:"id"`
			}{}
			t.Log(state.Trace, state.Version)
			//json.NewDecoder(context.Request.Body).Decode(&id)

			err = state.ShouldBindJSON(&id)
			// ValidationErrors

			assert.Nil(t, err)
			body, _ := state.Gin.Get(gin.BodyBytesKey)
			t.Log(body)
			t.Log(id)
			state.Success("ok")
		}))
		err = http.Serve(l, _engine)
		assert.Nil(t, err)
	}()

	time.Sleep(time.Second)
	var id = struct {
		Id int `json:"id"`
	}{Id: 10}
	var buf bytes.Buffer
	assert.Nil(t, json.NewEncoder(&buf).Encode(id))
	_, err = http.Post("http://localhost:8080/test?_trace=10313&_version=10.30", "Application/json", &buf)
	assert.Nil(t, err)
	_, err = http.Post("http://localhost:8080/1test?_trace=10314&_version=10.30&_uid=100", "Application/json", nil)
	assert.Nil(t, err)
}

func TestGroup(t *testing.T) {
	logger.SetLevel(logger.DebugLevel)
	viper.Set("Metrics.Address", ":18091")
	engine := Engine()
	engine.Use(func(context *gin.Context) {
		fmt.Println(">>global middleware")
		context.Next()
		fmt.Println("<<global middleware")
	})
	engine.Use(func(context *gin.Context) {
		fmt.Println("group middleware")
	}).GET("/ping", func(context *gin.Context) {
		fmt.Println("ping")
	})
	//i := 0
	engine.GET("/test", func(context *gin.Context) {
		state := NewState(context)
		state.httpJson(0, map[string]string{})
		//fmt.Println(2 / i)
	})
	neon.LoadPlugins(viper.GetViper())
	go func() {
		l, err := net.Listen("tcp", ":18090")

		if err != nil {
			logger.Errorf("listen %s fail, %v\n", ":18090", err)
		} else {
			logger.Infof("listen %s", l.Addr())
		}

		http.Serve(l, _engine)

	}()
	resp, _ := http.Get("http://127.0.0.1:18090/api/test")
	if resp != nil {
		buf, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(buf))
	}
	time.Sleep(time.Second * 1000)
}
