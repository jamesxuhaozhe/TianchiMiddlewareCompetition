package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/conf"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/constants"
	backend "github.com/jamesxuhaozhe/tianchimiddlewarecompetition/handler/backendprocess/engine"
	client "github.com/jamesxuhaozhe/tianchimiddlewarecompetition/handler/clientprocess/engine"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/log"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/router"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/utils"
	"github.com/spf13/pflag"
	"net/http"
	"time"
)

var (
	port = pflag.StringP("port", "p", "8080", "server port")
	mode = pflag.StringP("mode", "m", "release", "server mode: debug or release")
)

func main() {
	pflag.Parse()

	// SetServerPort the conf
	conf.SetServerPort(*port)

	// Start the logger
	log.InitLogger()

	// Set gin mode.
	gin.SetMode(*mode)

	// create gin engine
	g := gin.New()

	// routes
	router.Load(g)

	// client or backend process data structure init
	if utils.IsBackendProcess() {
		backend.Start()
		// init checksum goroutine
	} else if utils.IsClientProcess() {
		client.Start()
	}

	// Ping the server to make sure the router is working
	go func() {
		if err := pingServer(); err != nil {
			log.Fatal("The router has no response, or it might took too long to start up.")
		}
		log.Info("The router has been deployed successfully.")
	}()

	log.Infof("Start to listening the incoming requests on http address: %s",
		constants.CommonUrlPrefix+conf.GetServerPort())
	log.Fatal(http.ListenAndServe(":"+conf.GetServerPort(), g).Error())
}

func pingServer() error {

	for i := 0; i < constants.PingCount; i++ {
		resp, err := http.Get("http://" + constants.CommonUrlPrefix + conf.GetServerPort() + "/check/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}
		log.Info("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	return errors.New("cannot connect to the router")
}
