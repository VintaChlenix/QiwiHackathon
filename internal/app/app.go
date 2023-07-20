package app

type App struct {
	config *config.Config
}

func NewApp(config *config.Config) (*App, error) {
	return &App{
		config: config,
	}, nil
}
