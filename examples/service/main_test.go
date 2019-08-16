package service

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/geekymedic/neon/logger"
	"github.com/geekymedic/neon/plugin/rpc"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/geekymedic/neon/bff"
	"github.com/geekymedic/neon/service"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
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
	callhttp(t, ctx)
	cancel()
	time.Sleep(time.Second * 3)
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
		t.Log("trace", state.Trace, "version", state.Version)
		err = state.Gin.ShouldBindBodyWith(&id, binding.JSON)
		assert.Nil(t, err)
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
	resp, err := client.Ping(state.GrpcClientCtx(), &PingRequest{Msg:"http_client_request"})
	assert.Nil(t, err)
	t.Log(resp)
}

func callhttp(t *testing.T, ctx context.Context) {
	var id = struct {
		Id string `json:"id"`
	}{Id: uuid.Must(uuid.NewUUID()).String()}
	var buf bytes.Buffer
	assert.Nil(t, json.NewEncoder(&buf).Encode(id))
	_, err := http.Post("http://localhost:8080/api/id?_trace=10313&_version=10.30", "Application/json", &buf)
	assert.Nil(t, err)
}
