// @title Task Manager API
// @version 1.0
// @description Task management system with JWT authentication
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/maxzhirnov/go-task-manager/config"
	_ "github.com/maxzhirnov/go-task-manager/docs"
	"github.com/maxzhirnov/go-task-manager/internal/handlers"
	"github.com/maxzhirnov/go-task-manager/internal/middleware"
	"github.com/maxzhirnov/go-task-manager/pkg/database"
	"github.com/maxzhirnov/go-task-manager/pkg/email"
	httpSwagger "github.com/swaggo/http-swagger"
)

func setupRouter(cfg *config.Config) *mux.Router {
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize email service
	emailService := email.NewEmailService(
		cfg.SMTP.Host,
		cfg.SMTP.Port,
		cfg.SMTP.Username,
		cfg.SMTP.Password,
	)

	r := mux.NewRouter()

	// Auth handlers
	authHandler := handlers.NewAuthHandler(db, emailService)
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
	api.HandleFunc("/tasks/positions", taskHandler.UpdateTaskPositions).Methods("PUT")
	api.HandleFunc("/tasks/{id}", taskHandler.UpdateTask).Methods("PUT")
	api.HandleFunc("/tasks/{id}", taskHandler.DeleteTask).Methods("DELETE")

	api.HandleFunc("/users/statistics", taskHandler.GetUserStatistics).Methods("GET")

	// Swagger documentation
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
	))

	// Static files for Svelte assets (CSS, JS)
	r.PathPrefix("/_app/").Handler(http.FileServer(http.Dir("./frontend/build")))
	r.PathPrefix("/assets/").Handler(http.FileServer(http.Dir("./frontend/build")))

	// SPA fallback - must be last
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/build/index.html")
	})

	return r
}

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	r := setupRouter(cfg)

	serverAddr := ":" + cfg.Server.Port
	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(serverAddr, r))
}
