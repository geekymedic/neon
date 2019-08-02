package extend

import "testing"

func TestDefGrpLog(t *testing.T) {
	log := DefGrpLog()
	log.Info("Hello world")
	log.With("actor", "foo")
	log.Info("Continue")
}
