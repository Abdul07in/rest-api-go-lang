package domain

import (
	"context"
	"time"
)

type Student struct {
	ID        uint      `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	Grade     float64   `json:"grade"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type StudentRepository interface {
	Create(student *Student) error
	GetByID(id uint) (*Student, error)
	GetAll() ([]Student, error)
	Update(student *Student) error
	Delete(id uint) error
}

type StudentService interface {
	CreateStudent(ctx context.Context, student *Student) error
	GetStudent(ctx context.Context, id uint) (*Student, error)
	GetAllStudents(ctx context.Context) ([]Student, error)
	UpdateStudent(ctx context.Context, student *Student) error
	DeleteStudent(ctx context.Context, id uint) error
}
