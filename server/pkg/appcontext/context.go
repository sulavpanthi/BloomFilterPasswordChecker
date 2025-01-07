package appcontext

import (
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/sulavpanthi/BloomFilterPasswordChecker/pkg/config"
)

type AppContext struct {
	Config *config.Config
	Logger *zerolog.Logger
}

var (
	instance *AppContext
	once     sync.Once
)

func Initialize() error {
	var err error
	once.Do(func() {
		instance, err = setup()
	})
	return err
}

func setup() (*AppContext, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	// zerolog basic setup
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	log := zerolog.New(output).
		With().
		Timestamp().
		Caller().
		Logger()

	ctx := &AppContext{
		Config: cfg,
		Logger: &log,
	}

	ctx.Logger.Info().Msg("Application context initialized...")

	return ctx, nil
}

func Get() *AppContext {
	if instance == nil {
		panic("AppContext not initialized")
	}
	return instance
}

func Reset() {
	instance = nil
}
