package service

var beforeAppRun []func() error

var beforeAppExit []func() error

func RegisterBeforeAppRunHook(opts ...func() error) {
	for _, opt := range opts {
		beforeAppRun = append(beforeAppRun, opt)
	}
}

func RegisterBeforeAppExitHook(opts ...func() error) {
	for _, opt := range opts {
		beforeAppExit = append(beforeAppExit, opt)
	}
}
