package service

import (
	"errors"

	"gorm.io/gorm"

	"prince-university/student-records/model"
	"prince-university/student-records/repository"
)

var ErrCourseNotFound = errors.New("course not found")

type EnrollInput struct {
	MatricNumber    string
	FirstName       string
	LastName        string
	DateOfAdmission *string
	CourseID        uint
}

type StudentService struct {
	repo *repository.StudentRepository
}

func NewStudentService(repo *repository.StudentRepository) *StudentService {
	return &StudentService{repo: repo}
}

func (s *StudentService) ListStudents() ([]model.StudentView, error) {
	students, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	return toViews(students), nil
}

func (s *StudentService) ListHonorRoll() ([]model.StudentView, error) {
	students, err := s.repo.FindHonorRoll()
	if err != nil {
		return nil, err
	}
	return toViews(students), nil
}

func (s *StudentService) Enroll(in EnrollInput) (model.StudentView, error) {
	course, err := s.repo.FindCourse(in.CourseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.StudentView{}, ErrCourseNotFound
		}
		return model.StudentView{}, err
	}

	student := model.Student{
		MatricNumber: in.MatricNumber,
		FirstName:    in.FirstName,
		LastName:     in.LastName,
		Courses:      []model.Course{course},
	}
	if in.DateOfAdmission != nil && *in.DateOfAdmission != "" {
		d, err := model.ParseDate(*in.DateOfAdmission)
		if err != nil {
			return model.StudentView{}, errors.New("dateOfAdmission must be YYYY-MM-DD")
		}
		student.DateOfAdmission = &d
	}

	if err := s.repo.Create(&student); err != nil {
		return model.StudentView{}, err
	}
	return model.ToView(student), nil
}

func toViews(students []model.Student) []model.StudentView {
	views := make([]model.StudentView, 0, len(students))
	for _, s := range students {
		views = append(views, model.ToView(s))
	}
	return views
}
