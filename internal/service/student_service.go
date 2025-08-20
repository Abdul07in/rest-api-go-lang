package service

import "student-api/internal/domain"

type studentService struct {
	repo domain.StudentRepository
}

func NewStudentService(repo domain.StudentRepository) domain.StudentService {
	return &studentService{repo: repo}
}

func (s *studentService) CreateStudent(student *domain.Student) error {
	return s.repo.Create(student)
}

func (s *studentService) GetStudent(id uint) (*domain.Student, error) {
	return s.repo.GetByID(id)
}

func (s *studentService) GetAllStudents() ([]domain.Student, error) {
	return s.repo.GetAll()
}

func (s *studentService) UpdateStudent(student *domain.Student) error {
	return s.repo.Update(student)
}

func (s *studentService) DeleteStudent(id uint) error {
	return s.repo.Delete(id)
}
