package service

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/geekymedic/neon/logger"
	"github.com/geekymedic/neon/plugin/rpc"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"github.com/geekymedic/neon/bff"
	"github.com/geekymedic/neon/service"
)

type PingServer struct{}

func (s *PingServer) Ping(ctx context.Context, req *PingRequest) (*PingResponse, error) {
	time.Sleep(time.Second)
	return &PingResponse{}, nil
}

func TestPingCallChain(t *testing.T) {
	os.Setenv("NEON_MODE", "test")
	logger.SetLevel(logger.DebugLevel)
	ctx, cancel := context.WithCancel(context.TODO())
	go serviceServer(t, ctx)
	go bffServer(t, ctx)
	time.Sleep(time.Second)
	callhttpOk(t, ctx)
	// time.Sleep(time.Second)
	// callhttpFail(t, ctx)
	cancel()

}

func serviceServer(t *testing.T, ctx context.Context) {
	viper.Set("Address", ":8910")
	RegisterCheckHealthServer(service.Server(), &PingServer{})
	err := service.Main()
	assert.Nil(t, err)
}

func bffServer(t *testing.T, ctx context.Context) {
	l, err := net.Listen("tcp", ":8080")
	assert.Nil(t, err)

	g := bff.Engine()
	g.POST("/id", bff.HttpHandler(func(state *bff.State) {
		var id = struct {
			Id string `json:"id"`
		}{}
		err := state.ShouldBindJSON(&id)
		if err != nil {
			state.Error(bff.CodeRequestBodyError, err)
			return
		}
		t.Log("trace", state.Trace, "version", state.Version, "id", id.Id)
		callRpc(t, state)
		state.Success("ok")
	}))
	err = http.Serve(l, bff.MockEngine())
	assert.Nil(t, err)
}

func callRpc(t *testing.T, state *bff.State) {
	conn, err := grpc.Dial("localhost:8910", grpc.WithInsecure(), grpc.WithUnaryInterceptor(rpc.MockGrpcClientLog()))
	assert.Nil(t, err)
	client := NewCheckHealthClient(conn)
	resp, err := client.Ping(state.GrpcClientCtx(), &PingRequest{Msg: "http_client_request"})
	assert.Nil(t, err)
	t.Log(resp)
}

func callhttpOk(t *testing.T, ctx context.Context) {
	var id = struct {
		Id string `json:"id"`
	}{Id: uuid.Must(uuid.NewUUID()).String()}
	var buf bytes.Buffer
	assert.Nil(t, json.NewEncoder(&buf).Encode(id))
	_, err := http.Post("http://localhost:8080/api/id?_trace=10313&_version=10.30", "Application/json", &buf)
	assert.Nil(t, err)
}

func callhttpFail(t *testing.T, ctx context.Context) {
	var id = struct {
		Id interface{} `json:"id"`
	}{Id: 100}
	var buf bytes.Buffer
	assert.Nil(t, json.NewEncoder(&buf).Encode(id))
	_, err := http.Post("http://localhost:8080/api/id?_trace=10313&_version=10.30", "Application/json", &buf)
	assert.Nil(t, err)
}
