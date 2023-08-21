package api

import (
	"ibdatabase/controllers"
	"ibdatabase/grpcserver"
	"ibdatabase/repositories"
	"ibdatabase/repositories/mysqldb"
	"ibdatabase/services"
	"log"
	"os"
	"strconv"
)

// start and running api services
type Application struct {
	Controller *controllers.Controllers
	GrpcServer *grpcserver.GrpcServer
}

func (app *Application) StartApp() {
	//create and hold a controller service
	app.Controller = &controllers.Controllers{}
	app.Controller.Port = os.Getenv("WEB_PORT")

	//create and hold a db repo
	mysqlDBRepo := &mysqldb.MysqlDBRepo{}
	mysqlDBRepo.User = os.Getenv("MYSQL_USER")
	mysqlDBRepo.Password = os.Getenv("MYSQL_PASSWORD")
	mysqlDBRepo.DefaultDB = os.Getenv("MYSQL_DEFAULT_DB")

	deploy, err := strconv.ParseBool(os.Getenv("DEPLOY_DOCKER"))
	if err != nil {
		log.Fatal("deploy error: ", err)
	}
	if deploy {
		mysqlDBRepo.Host = os.Getenv("MYSQL_DOCKER_HOST")
	} else {
		mysqlDBRepo.Host = os.Getenv("MYSQL_HOST")
	}

	mysqlDBRepo.Port = os.Getenv("MYSQL_PORT")

	//assign db repo to service
	var dbrepo repositories.DatabaseRepo = mysqlDBRepo
	app.Controller.Services = services.Service{GetDBRepo: dbrepo}

	//grpc init
	app.GrpcServer = &grpcserver.GrpcServer{}
	app.GrpcServer.Port = os.Getenv("GRPC_PORT")
	enableTls, err := strconv.ParseBool(os.Getenv("GRPC_TLS"))

	if err != nil {
		log.Fatal(err)
	}

	app.GrpcServer.EnableTLS = enableTls
	app.GrpcServer.CertFile = os.Getenv("GRPC_CERT_FILE")
	app.GrpcServer.CertKey = os.Getenv("GRPC_CERT_KEY")
	app.GrpcServer.Services = &services.Service{GetDBRepo: dbrepo}

	//start grpc
	go app.GrpcServer.StartGRPC()

	//register routers and start server
	app.Controller.RegisterRouters()
}
