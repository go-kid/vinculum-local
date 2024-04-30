package vinculum_local

import (
	"github.com/fsnotify/fsnotify"
	"github.com/go-kid/ioc/configure/loader"
	"github.com/go-kid/ioc/syslog"
	"github.com/go-kid/properties"
	"github.com/go-kid/vinculum"
	"gopkg.in/yaml.v3"
)

type spy struct {
	Logger     syslog.Logger `logger:""`
	configPath string
	ch         chan<- properties.Properties
	watcher    *fsnotify.Watcher
}

func NewSpy(configPath string) vinculum.Spy {
	return &spy{
		configPath: configPath,
	}
}

func (s *spy) RegisterChannel(ch chan<- properties.Properties) {
	s.ch = ch
}

func (s *spy) Init() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	s.watcher = watcher
	return nil
}

func (s *spy) Run() error {
	go func() {
		for {
			select {
			case event, ok := <-s.watcher.Events:
				if !ok {
					return
				}
				if !event.Has(fsnotify.Write) {
					continue
				}
				config, err := loader.NewFileLoader(s.configPath).LoadConfig()
				if err != nil {
					s.Logger.Panicf("read file '%s' error: %v", s.configPath, err)
				}
				p := properties.New()
				err = yaml.Unmarshal(config, &p)
				if err != nil {
					s.Logger.Panicf("unmarshal config error: %v\n%s", err, string(config))
				}
				s.ch <- p
			case err, ok := <-s.watcher.Errors:
				if !ok {
					return
				}
				s.Logger.Panicf("watch file error: %v", err)
			}
		}
	}()
	err := s.watcher.Add(s.configPath)
	if err != nil {
		return err
	}
	return nil
}

func (s *spy) Close() error {
	return s.watcher.Close()
}
