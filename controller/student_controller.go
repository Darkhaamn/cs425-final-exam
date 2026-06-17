package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"prince-university/student-records/service"
)

type StudentController struct {
	svc *service.StudentService
}

func NewStudentController(svc *service.StudentService) *StudentController {
	return &StudentController{svc: svc}
}

func (ctl *StudentController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/students", ctl.ListStudents)
	rg.GET("/students/honor-roll", ctl.ListHonorRoll)
	rg.POST("/students", ctl.Enroll)
}

func (ctl *StudentController) ListStudents(c *gin.Context) {
	views, err := ctl.svc.ListStudents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, views)
}

func (ctl *StudentController) ListHonorRoll(c *gin.Context) {
	views, err := ctl.svc.ListHonorRoll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, views)
}

type enrollRequest struct {
	MatricNumber    string  `json:"matricNumber" binding:"required"`
	FirstName       string  `json:"firstName" binding:"required"`
	LastName        string  `json:"lastName" binding:"required"`
	DateOfAdmission *string `json:"dateOfAdmission"`
	CourseID        uint    `json:"courseId" binding:"required"`
}

func (ctl *StudentController) Enroll(c *gin.Context) {
	var req enrollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	view, err := ctl.svc.Enroll(service.EnrollInput{
		MatricNumber:    req.MatricNumber,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		DateOfAdmission: req.DateOfAdmission,
		CourseID:        req.CourseID,
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrCourseNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, view)
}
