package types

// DefaultConfiguration Копирует конфигурацию по умолчанию и возвращает ссылку на новую копию.
func DefaultConfiguration() (ret *Configuration) {
	ret = &Configuration{
		Autostart:      defaultConfiguration.Autostart,
		RestartTimeout: defaultConfiguration.RestartTimeout,
		KillTimeout:    defaultConfiguration.KillTimeout,
	}

	return
}
