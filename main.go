package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/jandro-es/merlin/configs"
	"github.com/jandro-es/merlin/helpers"
	"github.com/jandro-es/merlin/middleware"
	"github.com/jandro-es/merlin/routes"
)

func main() {
	doneCh := make(chan bool, 1)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	router := mux.NewRouter()
	// configs.ConnectDB()
	configs.ParseConfigurations()

	go func() {
		routes.GovernanceRoutes(router)
		routes.ConfigurableRoutes(router)
		router.Use(middleware.HeadersValidator)
		router.Use(middleware.Subrequests)
		router.Use(middleware.ContentTypeApplicationJsonMiddleware)
		router.Use(middleware.PassthroughHeaders)
		err := http.ListenAndServe(fmt.Sprintf(":%d", 9090), router)
		helpers.ExitOnFail(err, "Failed to start HTTP server")
	}()

	go func() {
		sig := <-sigCh
		fmt.Printf("Received %s signal, exiting.\n", sig.String())
		doneCh <- true
	}()

	fmt.Printf("The application has started. Listening on port: %d\n", 9090)
	fmt.Println("Ctrl+C to exit.")

	<-doneCh

	fmt.Println("Adios!!")
}
