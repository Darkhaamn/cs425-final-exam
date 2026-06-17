package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"prince-university/student-records/model"
)

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func dsn() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		env("DB_HOST", "localhost"),
		env("DB_PORT", "5432"),
		env("DB_USER", "prince"),
		env("DB_PASSWORD", "prince"),
		env("DB_NAME", "student_records"),
	)
}

func Connect() *gorm.DB {
	var db *gorm.DB
	var err error

	for attempt := 1; attempt <= 10; attempt++ {
		db, err = gorm.Open(postgres.Open(dsn()), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("waiting for postgres (attempt %d/10): %v", attempt, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("connect postgres: %v", err)
	}

	if err := db.AutoMigrate(&model.Course{}, &model.Student{}); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	seed(db)
	return db
}

func mustDate(s string) *model.Date {
	d, err := model.ParseDate(s)
	if err != nil {
		log.Fatalf("bad date %q: %v", s, err)
	}
	return &d
}

func seed(db *gorm.DB) {
	var count int64
	db.Model(&model.Course{}).Count(&count)
	if count > 0 {
		return
	}

	courses := []model.Course{
		{CourseID: 1, CourseCode: "MTH5002", CourseName: "Pure Mathematics", CreditScore: 5},
		{CourseID: 2, CourseCode: "PHY2009", CourseName: "Applied Physics", CreditScore: 2},
		{CourseID: 3, CourseCode: "CS6011", CourseName: "Advanced Computing", CreditScore: 6},
	}
	if err := db.Create(&courses).Error; err != nil {
		log.Fatalf("seed courses: %v", err)
	}

	byID := map[uint]model.Course{}
	for _, c := range courses {
		byID[c.CourseID] = c
	}
	pick := func(ids ...uint) []model.Course {
		out := make([]model.Course, 0, len(ids))
		for _, id := range ids {
			out = append(out, byID[id])
		}
		return out
	}

	students := []model.Student{
		{StudentID: 1, MatricNumber: "E019", FirstName: "Jennifer", LastName: "White", DateOfAdmission: mustDate("2026-01-15"), Courses: pick(2, 1)},
		{StudentID: 2, MatricNumber: "B107", FirstName: "Ben", LastName: "Brown", DateOfAdmission: mustDate("2026-01-15"), Courses: pick(3)},
		{StudentID: 3, MatricNumber: "E724", FirstName: "Ali", LastName: "McCoist", DateOfAdmission: mustDate("2026-03-31"), Courses: pick(1, 3, 2)},
		{StudentID: 4, MatricNumber: "A771", FirstName: "Isaiah", LastName: "Washington", DateOfAdmission: mustDate("2026-01-17"), Courses: pick(2)},
	}
	if err := db.Create(&students).Error; err != nil {
		log.Fatalf("seed students: %v", err)
	}

	resetSeq(db, "students", "student_id")
	resetSeq(db, "courses", "course_id")

	log.Println("seeded initial University data")
}

func resetSeq(db *gorm.DB, table, col string) {
	sql := fmt.Sprintf(
		"SELECT setval(pg_get_serial_sequence('%s','%s'), COALESCE((SELECT MAX(%s) FROM %s), 1))",
		table, col, col, table,
	)
	if err := db.Exec(sql).Error; err != nil {
		log.Fatalf("reset seq %s: %v", table, err)
	}
}
