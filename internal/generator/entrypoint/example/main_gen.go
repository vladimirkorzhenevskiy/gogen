package main

import (
	"github.com/vladimirkorzhenevskiy/gogen/pkg/application"
	"github.com/vladimirkorzhenevskiy/gogen/pkg/log"

	leaderboardapicontroller "example/internal/controller/leaderboard/api"
	leaderboardpostgresrepository "example/internal/repository/leaderboard/postgres"
	leaderboardgetusecase "example/internal/usecase/leaderboard/get"
	postgresdriver "github.com/vladimirkorzhenevskiy/gogen/pkg/postgres"
)

type Config struct {
	Drivers struct {
		Postgres postgresdriver.Config `envconfig:"POSTGRES"`
	} `envconfig:"DRIVERS"`
	Controllers struct {
		LeaderboardApi leaderboardapicontroller.Config `envconfig:"API"`
	} `envconfig:"CONTROLLERS"`
}

type Dependencies struct {
	Drivers struct {
		Postgres *postgresdriver.Driver
	}
	Controllers struct {
		LeaderboardApi *leaderboardapicontroller.Controller
	}
	UseCases struct {
		LeaderboardGet *leaderboardgetusecase.UseCase
	}
	Repositories struct {
		LeaderboardPostgres *leaderboardpostgresrepository.Repository
	}
}

func main() {
	logger := log.New()
	tracer := otel.GetTracerProvider().Tracer("example")

	var (
		dependencies Dependencies
		cfg          Config
	)

	{
		driver, err := postgresdriver.New(cfg.Drivers.Postgres)
		if err != nil {
			logger.Fatal("failed to create new postgres driver", log.Error(err))
		}

		driver.WithLogger(logger)
		driver.WithTracer(tracer)

		dependencies.Drivers.Postgres = driver
	}
	{
		repository, err := leaderboardpostgresrepository.New()
		if err != nil {
			logger.Fatal("failed to create new postgres repository", log.Error(err))
		}

		repository.WithLogger(logger)
		repository.WithTracer(tracer)
		repository.WithDb(dependencies.Drivers.Postgres)

		dependencies.Repositories.LeaderboardPostgres = repository
	}
	{
		usecase, err := leaderboardgetusecase.New()
		if err != nil {
			logger.Fatal("failed to create new get usecase", log.Error(err))
		}

		usecase.WithLogger(logger)
		usecase.WithTracer(tracer)
		usecase.WithRepository(dependencies.Repositories.LeaderboardPostgres)

		dependencies.UseCases.LeaderboardGet = usecase
	}
	{
		controller, err := leaderboardapicontroller.New(cfg.Controllers.LeaderboardApi)
		if err != nil {
			logger.Fatal("failed to create new api controller", log.Error(err))
		}

		controller.WithLogger(logger)
		controller.WithTracer(tracer)
		controller.WithGetter(dependencies.UseCases.LeaderboardGet)

		dependencies.Controllers.LeaderboardApi = controller
	}

	app := application.New(
		dependencies.Controllers.LeaderboardApi,
	)

	if err := app.Run(); err != nil {
		logger.Fatal("failed to run application", log.Error(err))
	}
}
