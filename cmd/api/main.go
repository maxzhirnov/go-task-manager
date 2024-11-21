package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/maxzhirnov/go-task-manager/internal/handlers"
	"github.com/maxzhirnov/go-task-manager/internal/middleware"
	"github.com/maxzhirnov/go-task-manager/pkg/database"
)

func setupRouter() *mux.Router {
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	r := mux.NewRouter()

	// Auth handlers
	authHandler := handlers.NewAuthHandler(db)
	r.HandleFunc("/api/register", authHandler.RegisterHandler).Methods("POST")
	r.HandleFunc("/api/login", authHandler.LoginHandler).Methods("POST")
	r.HandleFunc("/api/refresh", authHandler.RefreshTokenHandler).Methods("POST")

	// Task handlers
	taskHandler := handlers.NewTaskHandler(db)
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.JWTAuthMiddleware)
	api.HandleFunc("/tasks", taskHandler.GetTasks).Methods("GET")
	api.HandleFunc("/tasks", taskHandler.CreateTask).Methods("POST")
	api.HandleFunc("/tasks/{id}", taskHandler.GetTask).Methods("GET")
	api.HandleFunc("/tasks/{id}", taskHandler.UpdateTask).Methods("PUT")
	api.HandleFunc("/tasks/{id}", taskHandler.DeleteTask).Methods("DELETE")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web")))

	return r
}

func main() {
	r := setupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
