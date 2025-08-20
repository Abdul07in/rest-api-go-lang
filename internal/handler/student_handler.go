package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"student-api/internal/domain"
	"student-api/internal/logging"

	"github.com/gorilla/mux"
)

type StudentHandler struct {
	service domain.StudentService
	logger  *logging.RequestLogger
}

func NewStudentHandler(service domain.StudentService, logger *logging.RequestLogger) *StudentHandler {
	return &StudentHandler{
		service: service,
		logger:  logger,
	}
}

func (h *StudentHandler) CreateStudent(w http.ResponseWriter, r *http.Request) {
	traceID := logging.GetTraceIDFromContext(r.Context())
	
	var student domain.Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		h.logger.LogOperation(traceID, "CreateStudent", fmt.Sprintf("Invalid request body: %v", err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.logger.LogOperation(traceID, "CreateStudent", fmt.Sprintf("Creating student: %s %s", student.FirstName, student.LastName))
	
	if err := h.service.CreateStudent(&student); err != nil {
		h.logger.LogOperation(traceID, "CreateStudent", fmt.Sprintf("Error creating student: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.LogOperation(traceID, "CreateStudent", fmt.Sprintf("Student created successfully with ID: %d", student.ID))
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(student)
}

func (h *StudentHandler) GetStudent(w http.ResponseWriter, r *http.Request) {
	traceID := logging.GetTraceIDFromContext(r.Context())
	
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.logger.LogOperation(traceID, "GetStudent", fmt.Sprintf("Invalid student ID: %s", vars["id"]))
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	h.logger.LogOperation(traceID, "GetStudent", fmt.Sprintf("Fetching student with ID: %d", id))
	
	student, err := h.service.GetStudent(uint(id))
	if err != nil {
		h.logger.LogOperation(traceID, "GetStudent", fmt.Sprintf("Error fetching student: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if student == nil {
		h.logger.LogOperation(traceID, "GetStudent", fmt.Sprintf("Student not found with ID: %d", id))
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	h.logger.LogOperation(traceID, "GetStudent", fmt.Sprintf("Successfully fetched student: %s %s", student.FirstName, student.LastName))
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

func (h *StudentHandler) GetAllStudents(w http.ResponseWriter, r *http.Request) {
	students, err := h.service.GetAllStudents()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

func (h *StudentHandler) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	var student domain.Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	student.ID = uint(id)

	if err := h.service.UpdateStudent(&student); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

func (h *StudentHandler) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteStudent(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
