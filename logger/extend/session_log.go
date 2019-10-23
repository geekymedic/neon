package extend

import (
	"github.com/geekymedic/neon"
	"github.com/geekymedic/neon/logger"
)

func NewSessionLog(session *neon.Session) logger.Logger {
	return logger.With("_uid", session.Uid,
		"_token", session.Token,
		"_platform", session.Platform,
		"_version", session.Version,
		"_net", session.Net,
		"_mobile", session.Mobile,
		"_os", session.OS,
		"_device", session.Device,
		"_describe", session.Describe,
		"_trace", session.Trace,
		"_sequence", session.Sequence,
		"_time", session.Time,
		"_path", session.Path)
}
