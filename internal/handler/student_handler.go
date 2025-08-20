package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"student-api/internal/domain"
	"student-api/internal/logging"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type ResponseChannel struct {
	Data  interface{}
	Error error
}

type StudentHandler struct {
	service domain.StudentService
	logger  *logging.RequestLogger
	workers int
	jobs    chan func()
	wg      sync.WaitGroup
}

func NewStudentHandler(service domain.StudentService, logger *logging.RequestLogger) *StudentHandler {
	h := &StudentHandler{
		service: service,
		logger:  logger,
		workers: 50, // Number of worker goroutines
		jobs:    make(chan func(), 100), // Buffer size of 100 jobs
	}

	// Start worker pool
	for i := 0; i < h.workers; i++ {
		h.wg.Add(1)
		go h.startWorker()
	}

	return h
}

func (h *StudentHandler) startWorker() {
	defer h.wg.Done()
	for job := range h.jobs {
		job()
	}
}

func (h *StudentHandler) scheduleJob(job func()) {
	h.jobs <- job
}

func (h *StudentHandler) CreateStudent(w http.ResponseWriter, r *http.Request) {
	traceID := logging.GetTraceIDFromContext(r.Context())
	respChan := make(chan ResponseChannel, 1)
	
	var student domain.Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		h.logger.LogOperation(traceID, "CreateStudent", fmt.Sprintf("Invalid request body: %v", err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.logger.LogOperation(traceID, "CreateStudent", fmt.Sprintf("Creating student: %s %s", student.FirstName, student.LastName))
	
	// Process asynchronously
	h.scheduleJob(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := h.service.CreateStudent(ctx, &student)
		respChan <- ResponseChannel{Data: student, Error: err}
	})

	// Wait for response with timeout
	select {
	case resp := <-respChan:
		if resp.Error != nil {
			h.logger.LogOperation(traceID, "CreateStudent", fmt.Sprintf("Error creating student: %v", resp.Error))
			http.Error(w, resp.Error.Error(), http.StatusInternalServerError)
			return
		}
		h.logger.LogOperation(traceID, "CreateStudent", fmt.Sprintf("Student created successfully with ID: %d", student.ID))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp.Data)
	case <-time.After(15 * time.Second):
		h.logger.LogOperation(traceID, "CreateStudent", "Operation timed out")
		http.Error(w, "Request timeout", http.StatusGatewayTimeout)
	}
}

func (h *StudentHandler) GetStudent(w http.ResponseWriter, r *http.Request) {
	traceID := logging.GetTraceIDFromContext(r.Context())
	respChan := make(chan ResponseChannel, 1)
	
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.logger.LogOperation(traceID, "GetStudent", fmt.Sprintf("Invalid student ID: %s", vars["id"]))
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	h.logger.LogOperation(traceID, "GetStudent", fmt.Sprintf("Fetching student with ID: %d", id))
	
	// Process asynchronously
	h.scheduleJob(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		student, err := h.service.GetStudent(ctx, uint(id))
		respChan <- ResponseChannel{Data: student, Error: err}
	})

	// Wait for response with timeout
	select {
	case resp := <-respChan:
		if resp.Error != nil {
			h.logger.LogOperation(traceID, "GetStudent", fmt.Sprintf("Error fetching student: %v", resp.Error))
			http.Error(w, resp.Error.Error(), http.StatusInternalServerError)
			return
		}

		student, ok := resp.Data.(*domain.Student)
		if !ok || student == nil {
			h.logger.LogOperation(traceID, "GetStudent", fmt.Sprintf("Student not found with ID: %d", id))
			http.Error(w, "Student not found", http.StatusNotFound)
			return
		}

		h.logger.LogOperation(traceID, "GetStudent", fmt.Sprintf("Successfully fetched student: %s %s", student.FirstName, student.LastName))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(student)
	case <-time.After(15 * time.Second):
		h.logger.LogOperation(traceID, "GetStudent", "Operation timed out")
		http.Error(w, "Request timeout", http.StatusGatewayTimeout)
	}
}

func (h *StudentHandler) GetAllStudents(w http.ResponseWriter, r *http.Request) {
	traceID := logging.GetTraceIDFromContext(r.Context())
	respChan := make(chan ResponseChannel, 1)

	h.logger.LogOperation(traceID, "GetAllStudents", "Fetching all students")

	// Process asynchronously
	h.scheduleJob(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		students, err := h.service.GetAllStudents(ctx)
		respChan <- ResponseChannel{Data: students, Error: err}
	})

	// Wait for response with timeout
	select {
	case resp := <-respChan:
		if resp.Error != nil {
			h.logger.LogOperation(traceID, "GetAllStudents", fmt.Sprintf("Error fetching students: %v", resp.Error))
			http.Error(w, resp.Error.Error(), http.StatusInternalServerError)
			return
		}

		students, ok := resp.Data.([]domain.Student)
		if !ok {
			h.logger.LogOperation(traceID, "GetAllStudents", "Error converting response data")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		h.logger.LogOperation(traceID, "GetAllStudents", fmt.Sprintf("Successfully fetched %d students", len(students)))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(students)
	case <-time.After(15 * time.Second):
		h.logger.LogOperation(traceID, "GetAllStudents", "Operation timed out")
		http.Error(w, "Request timeout", http.StatusGatewayTimeout)
	}
}

func (h *StudentHandler) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	traceID := logging.GetTraceIDFromContext(r.Context())
	respChan := make(chan ResponseChannel, 1)

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.logger.LogOperation(traceID, "UpdateStudent", fmt.Sprintf("Invalid student ID: %s", vars["id"]))
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	var student domain.Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		h.logger.LogOperation(traceID, "UpdateStudent", fmt.Sprintf("Invalid request body: %v", err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	student.ID = uint(id)

	h.logger.LogOperation(traceID, "UpdateStudent", fmt.Sprintf("Updating student with ID: %d", id))

	// Process asynchronously
	h.scheduleJob(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := h.service.UpdateStudent(ctx, &student)
		respChan <- ResponseChannel{Data: student, Error: err}
	})

	// Wait for response with timeout
	select {
	case resp := <-respChan:
		if resp.Error != nil {
			h.logger.LogOperation(traceID, "UpdateStudent", fmt.Sprintf("Error updating student: %v", resp.Error))
			http.Error(w, resp.Error.Error(), http.StatusInternalServerError)
			return
		}

		h.logger.LogOperation(traceID, "UpdateStudent", fmt.Sprintf("Successfully updated student: %s %s", student.FirstName, student.LastName))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp.Data)
	case <-time.After(15 * time.Second):
		h.logger.LogOperation(traceID, "UpdateStudent", "Operation timed out")
		http.Error(w, "Request timeout", http.StatusGatewayTimeout)
	}
}

func (h *StudentHandler) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	traceID := logging.GetTraceIDFromContext(r.Context())
	respChan := make(chan ResponseChannel, 1)

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.logger.LogOperation(traceID, "DeleteStudent", fmt.Sprintf("Invalid student ID: %s", vars["id"]))
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	h.logger.LogOperation(traceID, "DeleteStudent", fmt.Sprintf("Deleting student with ID: %d", id))

	// Process asynchronously
	h.scheduleJob(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := h.service.DeleteStudent(ctx, uint(id))
		respChan <- ResponseChannel{Error: err}
	})

	// Wait for response with timeout
	select {
	case resp := <-respChan:
		if resp.Error != nil {
			h.logger.LogOperation(traceID, "DeleteStudent", fmt.Sprintf("Error deleting student: %v", resp.Error))
			http.Error(w, resp.Error.Error(), http.StatusInternalServerError)
			return
		}

		h.logger.LogOperation(traceID, "DeleteStudent", fmt.Sprintf("Successfully deleted student with ID: %d", id))
		w.WriteHeader(http.StatusNoContent)
	case <-time.After(15 * time.Second):
		h.logger.LogOperation(traceID, "DeleteStudent", "Operation timed out")
		http.Error(w, "Request timeout", http.StatusGatewayTimeout)
	}
}
