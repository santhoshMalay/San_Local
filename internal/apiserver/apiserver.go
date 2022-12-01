package apiserver

import (
	"log"

	"context"

	_ "github.com/joho/godotenv/autoload"
	"github.com/zhuravlev-pe/course-watch/internal/config"
	"github.com/zhuravlev-pe/course-watch/internal/delivery/http"
	httpV1 "github.com/zhuravlev-pe/course-watch/internal/delivery/http/v1"
	"github.com/zhuravlev-pe/course-watch/internal/delivery/http/v1/auth"
	"github.com/zhuravlev-pe/course-watch/internal/repository"
	"github.com/zhuravlev-pe/course-watch/internal/repository/fake_repo"
	"github.com/zhuravlev-pe/course-watch/internal/server"
	"github.com/zhuravlev-pe/course-watch/internal/service"
	"github.com/zhuravlev-pe/course-watch/pkg/idgen"
	"github.com/zhuravlev-pe/course-watch/pkg/keygen"
	"github.com/zhuravlev-pe/course-watch/pkg/postgres"
	"github.com/zhuravlev-pe/course-watch/pkg/security"
)

// @title Course Watch API
// @version 1.0
// @description REST API for Course Watch App

// @host localhost:8080
// @BasePath /api/v1/

// @tag.name User
// @tag.description Managing user account

// @tag.name courses
// @tag.description Temporary endpoints for Swagger demo. To be removed

// @tag.name Authentication
// @tag.description Login, logout and other security related operations

// Run initializes whole application.
func Run() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	idGen, err := idgen.New(cfg.SnowflakeNode)
	if err != nil {
		log.Fatal(err)
	}

	pgConfig := postgres.NewPgConfig(
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Database,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pgClient, err := postgres.NewClient(ctx, pgConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer pgClient.Close()

	repos := &repository.Repositories{
		Courses: fake_repo.NewCourses(),
		Users:   repository.NewUsersRepo(pgClient),
	}

	bearerAuth, err := createAuthenticator(cfg)
	if err != nil {
		log.Fatal(err)
	}

	services := service.NewServices(service.Deps{
		Repos: repos,
		IdGen: idGen,
	})

	handler := http.NewHandler(services, bearerAuth)

	srv := server.NewServer(cfg, handler.Init())

	log.Print("Starting server")
	if err = srv.Run(); err != nil {
		log.Fatal(err)
	}
}

func createAuthenticator(cfg *config.Config) (httpV1.BearerAuthenticator, error) {
	// JwtHandler uses HMAC-SHA256 for signing, block size for SHA256 is 64 bytes, so the key size is the same
	key, err := keygen.Generate(cfg.JWTAuthentication.SigningKey, "bearer-auth.key", 64)
	if err != nil {
		return nil, err
	}
	jwtHandler := security.NewJwtHandler(
		cfg.JWTAuthentication.Issuer,
		cfg.JWTAuthentication.ExpectedAudience,
		cfg.JWTAuthentication.TargetAudience,
		cfg.JWTAuthentication.TokenTTL,
		key,
	)
	bearerAuth := auth.NewBearerAuthenticator(jwtHandler)
	return bearerAuth, nil
}
