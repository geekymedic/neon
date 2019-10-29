package bff

import (
	"github.com/gin-gonic/gin"
)

type BbfHandler func(state *State)

func HttpHandler(handler BbfHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		handler(NewState(ctx))
	}
}

type BbfHandlerFunc func(state *State) (interface{}, int, error)

func HttpHandleFunc(fn BbfHandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		state := NewState(ctx)
		data, code, err := fn(state)
		if code == CodeSuccess {
			state.Success(data)
		} else {
			state.Error(code, err)
		}
	}
}
