package main

import (
	"fmt"
	"log"
	"net/http"
	"student-api/internal/config"
	"student-api/internal/database"
	"student-api/internal/handler"
	"student-api/internal/logging"
	"student-api/internal/repository"
	"student-api/internal/service"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := database.Initialize(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize dependencies
	logger := logging.NewRequestLogger()
	studentRepo := repository.NewMySQLStudentRepository(db)
	studentService := service.NewStudentService(studentRepo)
	studentHandler := handler.NewStudentHandler(studentService, logger)

	// Set up router
	router := mux.NewRouter()

	// Add logging middleware
	router.Use(logger.LogRequest)
	
	// Student routes
	router.HandleFunc("/api/students", studentHandler.CreateStudent).Methods("POST")
	router.HandleFunc("/api/students", studentHandler.GetAllStudents).Methods("GET")
	router.HandleFunc("/api/students/{id:[0-9]+}", studentHandler.GetStudent).Methods("GET")
	router.HandleFunc("/api/students/{id:[0-9]+}", studentHandler.UpdateStudent).Methods("PUT")
	router.HandleFunc("/api/students/{id:[0-9]+}", studentHandler.DeleteStudent).Methods("DELETE")

	// Start server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
