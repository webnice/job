package types // import "gopkg.in/webnice/job.v1/types"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
//import ()

// DefaultConfiguration Копирует конфигурацию по умолчанию и возвращает ссылку на новую копию
func DefaultConfiguration() (ret *Configuration) {
	ret = &Configuration{
		Autostart:      defaultConfiguration.Autostart,
		RestartTimeout: defaultConfiguration.RestartTimeout,
		KillTimeout:    defaultConfiguration.KillTimeout,
	}

	return
}
