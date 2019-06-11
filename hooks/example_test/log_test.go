package example_test

import (
	log "github.com/sirupsen/logrus"
	"github.com/vvlgo/gotools/hooks"
	"testing"
)

type LogInfo struct {
	Action      string
	UserAccount string
	Result      string
}

func TestMyhooks(t *testing.T) {

	hooks.LogConf("dev", nil)
	info := LogInfo{}
	info.UserAccount = "test"
	info.Result = "ok"
	info.Action = "test"

	log.Info(hooks.ToMap(info))
	log.Error(hooks.ToMap(info))
	log.Warn(hooks.ToMap(info))
	log.Fatal(hooks.ToMap(info))
}
