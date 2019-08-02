package service

type flags struct {
	Debug bool
}

var (
	_flags = &flags{
		Debug: false,
	}
)

func IsDebug() bool {
	return _flags.Debug
}
