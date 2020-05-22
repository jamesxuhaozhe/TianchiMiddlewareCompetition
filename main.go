package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/conf"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/constants"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/router"
	"github.com/spf13/pflag"
	"log"
	"net/http"
	"time"
)

var (
	port = pflag.StringP("port", "p", "8080", "server port")
	mode = pflag.StringP("mode", "m", "debug", "server mode: debug or release")
)

func main() {
	pflag.Parse()

	// SetServerPort the conf
	conf.SetServerPort(*port)

	// Set gin mode.
	gin.SetMode(*mode)

	// create gin engine
	g := gin.New()

	// routes
	// TODO need to add more engine based on the process type
	router.Load(g)

	// Ping the server to make sure the router is working
	go func() {
		if err := pingServer(); err != nil {
			log.Fatal("The router has no response, or it might took too long to start up.", err)
		}
		log.Println("The router has been deployed successfully.")
	}()

	log.Printf("Start to listening the incoming requests on http address: %s\n",
		constants.CommonUrlPrefix+conf.GetServerPort())
	log.Fatal(http.ListenAndServe(":"+conf.GetServerPort(), g).Error())
}

func pingServer() error {

	for i := 0; i < constants.PingCount; i++ {
		resp, err := http.Get("http://" + constants.CommonUrlPrefix + conf.GetServerPort() + "/check/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}
		log.Println("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	return errors.New("cannot connect to the router")
}
