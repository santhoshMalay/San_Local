package apiserver

import (
	"github.com/zhuravlev-pe/course-watch/internal/delivery/http"
	"github.com/zhuravlev-pe/course-watch/internal/repository/mockrepo"
	"github.com/zhuravlev-pe/course-watch/internal/server"
	"github.com/zhuravlev-pe/course-watch/internal/service"
	"github.com/zhuravlev-pe/course-watch/pkg/idgen"
	"log"
)

// @title Course Watch API
// @version 1.0
// @description REST API for Course Watch App

// @host localhost:8080
// @BasePath /api/v1/

// @securityDefinitions.apikey AdminAuth
// @in header
// @name Authorization

// @securityDefinitions.apikey StudentsAuth
// @in header
// @name Authorization

// @securityDefinitions.apikey UsersAuth
// @in header
// @name Authorization

// Run initializes whole application.
func Run() {

	//TODO: read nodeId from config
	idGen, err := idgen.New(1)
	if err != nil {
		log.Fatal(err)
	}

	repos := mockrepo.New()
	services := service.NewServices(service.Deps{
		Repos: repos,
		IdGen: idGen,
	})
	handler := http.NewHandler(services)

	srv := server.NewServer(handler.Init())

	log.Print("Starting server")
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
