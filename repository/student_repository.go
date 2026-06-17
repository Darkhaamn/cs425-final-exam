package repository

import (
	"gorm.io/gorm"

	"prince-university/student-records/model"
)

type StudentRepository struct {
	db *gorm.DB
}

func NewStudentRepository(db *gorm.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

func (r *StudentRepository) FindAll() ([]model.Student, error) {
	var students []model.Student
	err := r.db.Preload("Courses").Order("last_name asc").Find(&students).Error
	return students, err
}

func (r *StudentRepository) FindHonorRoll() ([]model.Student, error) {
	var students []model.Student
	err := r.db.Preload("Courses").
		Joins("JOIN student_courses sc ON sc.student_student_id = students.student_id").
		Joins("JOIN courses ON courses.course_id = sc.course_course_id").
		Where("courses.credit_score >= ?", 5).
		Group("students.student_id").
		Order("students.last_name asc").
		Find(&students).Error
	return students, err
}

func (r *StudentRepository) FindCourse(courseID uint) (model.Course, error) {
	var course model.Course
	err := r.db.First(&course, courseID).Error
	return course, err
}

func (r *StudentRepository) Create(s *model.Student) error {
	return r.db.Create(s).Error
}
