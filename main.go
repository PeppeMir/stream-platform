package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"stream-platform/databases"
	"stream-platform/endpoints"
	"stream-platform/handlers"

	"github.com/gorilla/mux"
)

func main() {
	router := ConfigureRouter()

	databases.Connect()

	slog.Info("Listening on port " + os.Getenv("PORT"))
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), router)

}

func ConfigureRouter() *mux.Router {
	router := mux.NewRouter()

	usersRouter := router.PathPrefix("/api/users").Subrouter().StrictSlash(true)
	usersRouter.HandleFunc("/register", endpoints.CreateUser).Methods("POST")
	usersRouter.HandleFunc("/auth", endpoints.Authenticate).Methods("POST")

	moviesRouter := router.PathPrefix("/api/movies").Subrouter().StrictSlash(false)
	moviesRouter.Use(handlers.AuthHandler) // Authenticated requests via JWT
	moviesRouter.HandleFunc("/search", endpoints.SearchMovies).Methods("GET")
	moviesRouter.HandleFunc("/{id}", endpoints.GetMovie).Methods("GET")
	moviesRouter.HandleFunc("", endpoints.CreateMovie).Methods("POST")
	moviesRouter.HandleFunc("", endpoints.UpdateMovie).Methods("PUT")
	moviesRouter.HandleFunc("/{id}", endpoints.DeleteMovie).Methods("DELETE")

	return router
}
