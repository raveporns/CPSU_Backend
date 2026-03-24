package handler

import (
	"database/sql"
	"errors"
	"net/http"

	"cpsu/internal/course/models"
	"cpsu/internal/course/service"

	"github.com/gin-gonic/gin"
)

type CourseHandler struct {
	courseService service.CourseService
}

func NewCourseHandler(courseService service.CourseService) *CourseHandler {
	return &CourseHandler{courseService: courseService}
}

func (h *CourseHandler) GetAllCourses(c *gin.Context) {
	var param models.CoursesQueryParam
	if err := c.BindQuery(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	courses, err := h.courseService.GetAllCourses(param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, courses)
}

func (h *CourseHandler) GetCourseByID(c *gin.Context) {
	id := c.Param("id")

	course, err := h.courseService.GetCourseByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, course)
}

func (h *CourseHandler) CreateCourse(c *gin.Context) {
	var req models.CoursesRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetInt("user_id")
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	createdCourse, err := h.courseService.CreateCourse(req, userID, ip, userAgent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdCourse)
}

func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	id := c.Param("id")

	var req models.CoursesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetInt("user_id")
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	updatedCourse, err := h.courseService.UpdateCourse(id, req, userID, ip, userAgent)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, updatedCourse)
}

func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	id := c.Param("id")

	userID := c.GetInt("user_id")
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	err := h.courseService.DeleteCourse(id, userID, ip, userAgent)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "course deleted successfully"})
}
