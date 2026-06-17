package model

type Course struct {
	CourseID    uint      `gorm:"primaryKey;column:course_id" json:"courseId"`
	CourseCode  string    `gorm:"uniqueIndex;not null" json:"courseCode"`
	CourseName  string    `gorm:"not null" json:"courseName"`
	CreditScore int       `gorm:"not null" json:"creditScore"`
	Students    []Student `gorm:"many2many:student_courses;" json:"-"`
}

type Student struct {
	StudentID       uint     `gorm:"primaryKey;column:student_id" json:"studentId"`
	MatricNumber    string   `gorm:"uniqueIndex;not null" json:"matricNumber"`
	FirstName       string   `gorm:"not null" json:"firstName"`
	LastName        string   `gorm:"not null" json:"lastName"`
	DateOfAdmission *Date    `gorm:"type:date" json:"dateOfAdmission,omitempty"`
	Courses         []Course `gorm:"many2many:student_courses;" json:"courses"`
}

type StudentView struct {
	StudentID        uint     `json:"studentId"`
	MatricNumber     string   `json:"matricNumber"`
	FirstName        string   `json:"firstName"`
	LastName         string   `json:"lastName"`
	DateOfAdmission  *Date    `json:"dateOfAdmission,omitempty"`
	Courses          []Course `json:"courses"`
	TotalCreditScore int      `json:"totalCreditScore"`
}

func ToView(s Student) StudentView {
	total := 0
	for _, c := range s.Courses {
		total += c.CreditScore
	}
	return StudentView{
		StudentID:        s.StudentID,
		MatricNumber:     s.MatricNumber,
		FirstName:        s.FirstName,
		LastName:         s.LastName,
		DateOfAdmission:  s.DateOfAdmission,
		Courses:          s.Courses,
		TotalCreditScore: total,
	}
}
