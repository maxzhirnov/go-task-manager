package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/maxzhirnov/go-task-manager/internal/handlers"
	"github.com/maxzhirnov/go-task-manager/pkg/database"
)

func setupRouter() *mux.Router {
	// Initialize database connection
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create router
	r := mux.NewRouter()

	// Create handler
	taskHandler := handlers.NewTaskHandler(db)

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/tasks", taskHandler.GetTasks).Methods("GET")
	api.HandleFunc("/tasks", taskHandler.CreateTask).Methods("POST")
	api.HandleFunc("/tasks/{id}", taskHandler.GetTask).Methods("GET")
	api.HandleFunc("/tasks/{id}", taskHandler.UpdateTask).Methods("PUT")
	api.HandleFunc("/tasks/{id}", taskHandler.DeleteTask).Methods("DELETE")

	// Serve a default route, e.g., a homepage or a health check
	// r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusOK)
	// 	w.Write([]byte("Welcome to the Task Manager API"))
	// })

	// Serve static files
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web")))

	return r
}

func main() {
	// Get the router
	r := setupRouter()

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
