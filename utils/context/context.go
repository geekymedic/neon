package context

import (
	"context"

	"github.com/geekymedic/neon"
	"google.golang.org/grpc/metadata"
)

func CreateContextWithSession(ctx context.Context,x *neon.Session) context.Context {

	values := map[string]string{}
	values["_uid"] = x.Uid
	values["_trace"] = x.Trace
	values["_token"] = x.Token
	values["_describe"] = x.Describe
	values["_device"] = x.Device
	values["_mobile"] = x.Mobile
	values["_net"] = x.Net
	values["_os"] = x.OS
	values["_platform"] = x.Platform
	values["_sequence"] = x.Sequence
	values["_time"] = x.Time
	values["_version"] = x.Version
	md := metadata.New(values)


	return metadata.NewOutgoingContext(ctx, md)
}
