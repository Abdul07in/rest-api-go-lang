package service

import (
	"context"
	"student-api/internal/domain"
)

type studentService struct {
	repo domain.StudentRepository
}

func NewStudentService(repo domain.StudentRepository) domain.StudentService {
	return &studentService{repo: repo}
}

func (s *studentService) CreateStudent(ctx context.Context, student *domain.Student) error {
	return s.repo.Create(student)
}

func (s *studentService) GetStudent(ctx context.Context, id uint) (*domain.Student, error) {
	return s.repo.GetByID(id)
}

func (s *studentService) GetAllStudents(ctx context.Context) ([]domain.Student, error) {
	return s.repo.GetAll()
}

func (s *studentService) UpdateStudent(ctx context.Context, student *domain.Student) error {
	return s.repo.Update(student)
}

func (s *studentService) DeleteStudent(ctx context.Context, id uint) error {
	return s.repo.Delete(id)
}
