package routes

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jandro-es/merlin/controllers"
)

func LoginRoutes(router *mux.Router) {
	router.Handle("/login", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(controllers.Login()))).Methods("POST")
}
