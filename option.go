package vinculum_local

import "github.com/go-kid/ioc/app"

func Plugin(configPath string) app.SettingOption {
	return app.Options(
		app.SetConfig(configPath),
		app.SetComponents(NewSpy(configPath)),
	)
}
