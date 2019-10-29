package bff

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/geekymedic/neon/bff/types"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"

	"github.com/geekymedic/neon"
	"github.com/geekymedic/neon/errors"
	"github.com/geekymedic/neon/logger"
	"github.com/geekymedic/neon/logger/extend"
)

var (
	empty = struct{}{}
)

func newSessionCtx(ctx context.Context, session *neon.Session) context.Context {

	return metadata.NewOutgoingContext(
		ctx,
		metadata.New(session.KeysValues()),
	)
}

func NewState(ctx *gin.Context) *State {

	x, exists := ctx.Get(types.StateName)

	if exists {
		return x.(*State)
	}

	var (
		context = context.Background()
		session = neon.NewSessionFromGinCtx(ctx)
		state   = &State{
			Gin:     ctx,
			Session: session,
			Logger:  extend.NewSessionLog(session),
			ctx:     newSessionCtx(context, session),
		}
	)

	ctx.Set(types.StateName, state)

	return state
}

type State struct {
	*neon.Session
	logger.Logger
	Gin *gin.Context
	ctx context.Context
}

func (m *State) Error(code int, err error) {
	m.Gin.Set(types.ResponseStatusCode, code)
	if err != nil {
		m.Gin.Set(types.ResponseErr, err)
	}
	m.httpJson(code, empty)
}

func (m *State) ErrorMessage(code int, txt string) {
	m.Gin.Set(types.ResponseStatusCode, code)
	if txt != "" {
		m.Gin.Set(types.ResponseErr, fmt.Errorf("%s", txt))
	}
	m.httpJsonMessage(code, txt, empty)
}

func (m *State) Success(v interface{}) {
	m.Gin.Set(types.ResponseStatusCode, CodeSuccess)
	m.httpJson(
		CodeSuccess,
		v)
}

func (m *State) Context() context.Context {
	return m.ctx
}

func (m *State) GrpcClientCtx() context.Context {
	return metadata.NewOutgoingContext(context.Background(), metadata.New(m.Session.KeysValues()))
}

func (m *State) ShouldBindJSON(v interface{}) error {
	if m.Gin.Writer.Written() {
		return nil
	}
	if err := m.Gin.ShouldBindJSON(v); err != nil {
		m.Session.StructError = err.Error()
		m.Gin.Set(types.NeonSession, m.Session)
		return errors.By(err)
	}
	return nil
}

func (m *State) httpJson(code int, v interface{}) {
	if m.Gin.Writer.Written() {
		return
	}
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(map[string]interface{}{
		"Code":    code,
		"Message": GetMessage(code),
		"Data":    v,
	})
	if err != nil {
		logger.With("err", err).Error("fail to encode response")
	}
	m.Gin.Writer.WriteHeader(http.StatusOK)
	_, err = m.Gin.Writer.Write(buf.Bytes())
	if err != nil {
		m.Gin.Set(types.ResponseErr, err)
	}
}

func (m *State) httpJsonMessage(code int, message string, v interface{}) {
	if m.Gin.Writer.Written() {
		return
	}
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(map[string]interface{}{
		"Code":    code,
		"Message": message,
		"Data":    v,
	})
	if err != nil {
		logger.With("err", err).Error("fail to encode response")
	}
	m.Gin.Writer.WriteHeader(http.StatusOK)
	_, err = m.Gin.Writer.Write(buf.Bytes())
	if err != nil {
		m.Gin.Set(types.ResponseErr, err)
	}
}
