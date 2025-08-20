package repository

import (
	"database/sql"
	"student-api/internal/domain"
	"time"
)

type mysqlStudentRepository struct {
	db *sql.DB
}

func NewMySQLStudentRepository(db *sql.DB) domain.StudentRepository {
	return &mysqlStudentRepository{db: db}
}

func (r *mysqlStudentRepository) Create(student *domain.Student) error {
	query := `
		INSERT INTO students (first_name, last_name, email, age, grade, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	result, err := r.db.Exec(query,
		student.FirstName,
		student.LastName,
		student.Email,
		student.Age,
		student.Grade,
		now,
		now,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	student.ID = uint(id)
	student.CreatedAt = now
	student.UpdatedAt = now
	return nil
}

func (r *mysqlStudentRepository) GetByID(id uint) (*domain.Student, error) {
	query := `
		SELECT id, first_name, last_name, email, age, grade, created_at, updated_at
		FROM students
		WHERE id = ?
	`
	student := &domain.Student{}
	err := r.db.QueryRow(query, id).Scan(
		&student.ID,
		&student.FirstName,
		&student.LastName,
		&student.Email,
		&student.Age,
		&student.Grade,
		&student.CreatedAt,
		&student.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return student, nil
}

func (r *mysqlStudentRepository) GetAll() ([]domain.Student, error) {
	query := `
		SELECT id, first_name, last_name, email, age, grade, created_at, updated_at
		FROM students
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []domain.Student
	for rows.Next() {
		var student domain.Student
		err := rows.Scan(
			&student.ID,
			&student.FirstName,
			&student.LastName,
			&student.Email,
			&student.Age,
			&student.Grade,
			&student.CreatedAt,
			&student.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}
	return students, nil
}

func (r *mysqlStudentRepository) Update(student *domain.Student) error {
	query := `
		UPDATE students
		SET first_name = ?, last_name = ?, email = ?, age = ?, grade = ?, updated_at = ?
		WHERE id = ?
	`
	now := time.Now()
	_, err := r.db.Exec(query,
		student.FirstName,
		student.LastName,
		student.Email,
		student.Age,
		student.Grade,
		now,
		student.ID,
	)
	if err != nil {
		return err
	}
	student.UpdatedAt = now
	return nil
}

func (r *mysqlStudentRepository) Delete(id uint) error {
	query := "DELETE FROM students WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}
