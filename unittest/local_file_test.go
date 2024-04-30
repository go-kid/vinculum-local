package unittest

import (
	"fmt"
	"github.com/go-kid/ioc"
	"github.com/go-kid/ioc/app"
	"github.com/go-kid/vinculum"
	vinculum_local "github.com/go-kid/vinculum-local"
	"testing"
)

type Config struct {
	F1 string `yaml:"f1"`
}

func (c *Config) Prefix() string {
	return "Config"
}

type TApp struct {
	Config *Config `refreshScope:""`
}

func (t *TApp) Init() error {
	fmt.Printf("TApp inited\n%+v\n", t.Config)
	return nil
}

func (t *TApp) OnScopeChange(path string) error {
	fmt.Printf("refresh on %s\n%+v\n", path, t.Config)
	return nil
}

func TestLocalConfigSpy(t *testing.T) {
	ioc.RunTest(t,
		app.LogTrace,
		vinculum_local.Plugin("./config.yaml"),
		app.SetComponents(vinculum.New(), &TApp{}),
	)
	<-make(chan struct{})
}
