package service

import (
	"context"
	"github.com/geekymedic/neon/logger/extend"
	"time"

	"github.com/geekymedic/neon/logger"

	"github.com/geekymedic/neon"
	"google.golang.org/grpc/metadata"
)

var (
	stateName = "neon.service.State"
)

type State struct {
	logger.Logger
	ctx context.Context
	*neon.Session
}

func (m *State) Context() context.Context {
	return m.ctx
}

func (m *State) ContextWithTimeout(timeout time.Duration) context.Context {
	ctx, _ := context.WithTimeout(m.ctx, timeout)
	return ctx
}

// GrpcClientCtx return a new grpc client context
func (m *State) GrpcClientCtx() context.Context {
	return metadata.NewOutgoingContext(context.Background(), metadata.New(m.Session.KeysValues()))
}

func NewState(ctx context.Context) *State {

	v := ctx.Value(stateName)

	state, ok := v.(*State)

	if ok {
		return state
	}

	var (
		session    = &neon.Session{}
		md, exists = metadata.FromIncomingContext(ctx)
		value      = func(name string, x *string) {

			data := md.Get(name)

			if len(data) > 0 {
				*x = data[0]
			}

		}
	)

	if exists {
		for name, ref := range session.Keys() {
			value(name, ref)
		}
	}

	ctx = metadata.NewOutgoingContext(
		ctx,
		metadata.New(
			session.KeysValues(),
		),
	)

	state = &State{
		Session: session,
		Logger:  extend.NewSessionLog(session),
	}
	ctx = context.WithValue(ctx, stateName, state)

	state.ctx = ctx

	return state

}
