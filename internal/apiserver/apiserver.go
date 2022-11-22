package apiserver

import (
<<<<<<< HEAD
	"log"

=======
	"github.com/zhuravlev-pe/course-watch/internal/config"
>>>>>>> 64b2afd3fc6c83fbffcb07e6c3eb3f3d3bca3d1a
	"github.com/zhuravlev-pe/course-watch/internal/delivery/http"
	"github.com/zhuravlev-pe/course-watch/internal/delivery/http/v1/fake_authenticator"
	"github.com/zhuravlev-pe/course-watch/internal/repository/fake_repo"
	"github.com/zhuravlev-pe/course-watch/internal/server"
	"github.com/zhuravlev-pe/course-watch/internal/service"
	"github.com/zhuravlev-pe/course-watch/pkg/idgen"
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
	
	repos := fake_repo.New()
	
	//jwtHandler := security.NewJwtHandler([]byte(cfg.SigningKey))
	// TODO: first config related task - configure jwtHandler
	
	fakeBearerAuth := fake_authenticator.New() // to be able to test /user endpoints without logging in
	
	services := service.NewServices(service.Deps{
		Repos: repos,
		IdGen: idGen,
	})
	handler := http.NewHandler(services, fakeBearerAuth)
	
	srv := server.NewServer(cfg, handler.Init())
	
	log.Print("Starting server")
	if err = srv.Run(); err != nil {
		log.Fatal(err)
	}
}
