package main

import (
	"context"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/louispy/miniloan/internal/api"
	appAPI "github.com/louispy/miniloan/internal/api"
	"github.com/louispy/miniloan/internal/database"
	"github.com/louispy/miniloan/internal/domain/repositories"
	"github.com/louispy/miniloan/internal/services"
)

type Container struct {
	API *api.API
}

type Config struct {
	DB                    database.Config `envPrefix:"DB_"`
	BusinessTZOffsetHours int             `env:"BUSINESS_TZ_OFFSET_HOURS"`
}

func NewConfig() *Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		panic("cannot parse config: " + err.Error())
	}

	return &cfg
}
func NewContainer() *Container {
	cfg := NewConfig()

	db, err := database.New(context.Background(), cfg.DB)
	if err != nil {
		panic("database cannot be initialized: " + err.Error())
	}
	loanRepo := repositories.NewLoansRepository(repositories.LoanRepoOpts{DB: db})
	installmentRepo := repositories.NewInstallmentsRepository(repositories.InstallmentRepoOpts{DB: db})
	txManager := database.NewTxManager(database.TxManagerOpts{DB: db})
	businessTZ := time.FixedZone("BIZ", cfg.BusinessTZOffsetHours*3600)
	loanService := services.NewLoanService(services.LoanServiceOpts{
		LoansRepo:        loanRepo,
		InstallmentsRepo: installmentRepo,
		TxManager:        txManager,
		BusinessTZ:       businessTZ,
	})
	api := appAPI.NewAPI(appAPI.Opts{
		LoanService: loanService,
	})
	api.Register()

	return &Container{
		API: api,
	}
}
