package routes

import (
	"github.com/gorilla/mux"
	"github.com/jandro-es/merlin/controllers"
)

func GovernanceRoutes(router *mux.Router) {
	router.HandleFunc("/liveness", controllers.Liveness()).Methods("GET")
	router.HandleFunc("/readiness", controllers.Readiness()).Methods("GET")
}
