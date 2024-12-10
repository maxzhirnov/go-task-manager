package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/maxzhirnov/go-task-manager/internal/handlers"
	"github.com/maxzhirnov/go-task-manager/internal/middleware"
	"github.com/maxzhirnov/go-task-manager/pkg/config"
	"github.com/maxzhirnov/go-task-manager/pkg/database"
	"github.com/maxzhirnov/go-task-manager/pkg/email"
)

func setupRouter(cfg *config.Config) *mux.Router {
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize email service
	emailService, err := email.NewEmailService(
		cfg.SMTP.Host,
		cfg.SMTP.Port,
		cfg.SMTP.Username,
		cfg.SMTP.From,
		cfg.SMTP.FromName,
		cfg.SMTP.Password,
		cfg.SMTP.BaseURL,
	)
	if err != nil {
		log.Fatalf("Failed to initialize email service: %v", err)
	}

	r := mux.NewRouter()

	// Auth handlers
	authHandler := handlers.NewAuthHandler(db, emailService, cfg)
	r.HandleFunc("/api/register", authHandler.RegisterHandler).Methods("POST")
	r.HandleFunc("/api/login", authHandler.LoginHandler).Methods("POST")
	r.HandleFunc("/api/refresh", authHandler.RefreshTokenHandler).Methods("POST")
	r.HandleFunc("/api/verify-email", authHandler.VerifyEmailHandler).Methods("GET")
	r.HandleFunc("/api/resend-verification", authHandler.ResendVerificationHandler).Methods("POST")
	r.HandleFunc("/api/forgot-password", authHandler.ForgotPasswordHandler).Methods("POST")
	r.HandleFunc("/api/reset-password", authHandler.ResetPasswordHandler).Methods("POST")

	// Task handlers
	taskHandler := handlers.NewTaskHandler(db)
	userHandler := handlers.NewUserHandler(db)
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.JWTAuthMiddleware)
	api.HandleFunc("/tasks", taskHandler.GetTasks).Methods("GET")
	api.HandleFunc("/tasks", taskHandler.CreateTask).Methods("POST")
	api.HandleFunc("/tasks/{id}", taskHandler.GetTask).Methods("GET")
	api.HandleFunc("/tasks/positions", taskHandler.UpdateTaskPositions).Methods("PUT")
	api.HandleFunc("/tasks/{id}", taskHandler.UpdateTask).Methods("PUT")
	api.HandleFunc("/tasks/{id}", taskHandler.DeleteTask).Methods("DELETE")

	api.HandleFunc("/users/statistics", taskHandler.GetUserStatistics).Methods("GET")
	api.HandleFunc("/profile", userHandler.UpdateProfile).Methods("PUT")

	// Static files for Svelte assets (CSS, JS)
	r.PathPrefix("/_app/").Handler(http.FileServer(http.Dir("./frontend/build")))
	r.PathPrefix("/assets/").Handler(http.FileServer(http.Dir("./frontend/build")))

	r.HandleFunc("/favicon.svg", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		http.ServeFile(w, r, "./frontend/build/favicon.svg")
	})

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
