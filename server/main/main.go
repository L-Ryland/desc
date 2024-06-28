package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"server/gateway"
	"server/mongodb"
	"server/usersys"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

var (
	ifInit        = flag.Bool("init", false, "if init server")
	enableSwagger = flag.Bool("swagger", false, "if serve swagger")
	debugMode     = flag.Bool("debug", false, "if print debug log")
	mongo         = flag.Bool("mongodb", true, "use mongodb")
	flagPort      = flag.String("port", "8071", "server port. Default as 8081.")
)

var router = mux.NewRouter()

func init() {
	http.Handle("/", router)
}

func main() {
	flag.Parse()

	// logger
	logrus.SetReportCaller(true)
	customFormatter := new(logrus.TextFormatter)
	customFormatter.FullTimestamp = true
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	logrus.SetFormatter(customFormatter)
	if *debugMode {
		logrus.SetLevel(logrus.DebugLevel)
	}

	mongodb.NewDatabase()
	logrus.Info("database connected")

	usersys.Init()

	// network
	gateway.NewService(router)

	// 启动服务
	log.Println("Server started at http://localhost:" + *flagPort + "/")
	c := cors.New(cors.Options{
		// AllowedOrigins:   []string{"*"},
		AllowedOrigins:   []string{"http://localhost:8080", "http://localhost:3001"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPatch},
		AllowCredentials: true,
	})
	handler := c.Handler(router)
	if err := http.ListenAndServe(":"+*flagPort, handler); err != nil {
		log.Fatal(err)
	}
}
