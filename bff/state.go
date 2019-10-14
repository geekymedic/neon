package bff

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/geekymedic/neon/bff/types"

	"github.com/geekymedic/neon"
	"github.com/geekymedic/neon/errors"
	"github.com/geekymedic/neon/logger"
	"github.com/geekymedic/neon/logger/extend"
	"github.com/geekymedic/neon/version"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

var (
	empty = struct{}{}
)

type State struct {
	*neon.Session
	logger.Logger
	Gin *gin.Context
	ctx context.Context
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
	log := logger.With(m.Session.ShortLog()...).With("gitcommit", version.GITCOMMIT)
	if err != nil || code != CodeSuccess {
		log.With("body", buf.String()).With("err", err).Error("http response trace")
	} else {
		log.With("body", buf.String()).Info("http response trace")
	}

	m.Gin.Set(types.ResponseStatusCode, code)
}

func (m *State) httpJsonMessage(code int, message string, v interface{}) {
	if m.Gin.Writer.Written() {
		return
	}
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(map[string]interface{}{
		"Code":    code,
		"Message": GetMessage(code) + "|" + message,
		"Data":    v,
	})
	if err != nil {
		logger.With("err", err).Error("fail to encode response")
	}
	m.Gin.Writer.WriteHeader(http.StatusOK)
	_, err = m.Gin.Writer.Write(buf.Bytes())
	log := logger.With(m.Session.ShortLog()...).With("gitcommit", version.GITCOMMIT)
	if err != nil || code != CodeSuccess {
		log.With("body", buf.String()).With("err", err).Error("http response trace")
	} else {
		log.With("body", buf.String()).Info("http response trace")
	}

	m.Gin.Set(types.ResponseStatusCode, code)
}

func (m *State) Error(code int, err error) {
	if err != nil {
		m.Logger.Error(err.Error())
	}
	m.httpJson(code, empty)
}

func (m *State) ErrorMessage(code int, txt string) {
	m.Logger.Error(txt)
	m.httpJsonMessage(code, txt, empty)
}

func (m *State) Success(v interface{}) {
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

func (m *State) BindJSON(v interface{}) error {

	err := m.Gin.BindJSON(v)

	if err != nil {
		return errors.By(err)
	}

	return nil
}

func (m *State) ShouldBindJSON(v interface{}) error {
	if m.Gin.Writer.Written() {
		return nil
	}
	return errors.By(m.Gin.ShouldBindJSON(v))
}

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
