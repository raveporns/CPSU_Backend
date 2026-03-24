package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"cpsu/internal/course_structure/models"
	"cpsu/internal/course_structure/service"

	"github.com/gin-gonic/gin"
)

type CourseStructureHandler struct {
	courseStructureService service.CourseStructureService
}

func NewCourseStructureHandler(courseStructureService service.CourseStructureService) *CourseStructureHandler {
	return &CourseStructureHandler{courseStructureService: courseStructureService}
}

func (h *CourseStructureHandler) GetAllCourseStructure(c *gin.Context) {
	var param models.CourseStructureQueryParam
	if err := c.ShouldBindQuery(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	courseStructures, err := h.courseStructureService.GetAllCourseStructure(param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, courseStructures)
}

func (h *CourseStructureHandler) GetCourseStructureByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course_structure ID"})
		return
	}

	courseStructures, err := h.courseStructureService.GetCourseStructureByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "course_structure not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, courseStructures)
}

func (h *CourseStructureHandler) CreateCourseStructure(c *gin.Context) {

	courseID := c.PostForm("course_id")

	file, err := c.FormFile("detail")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer openedFile.Close()

	detail, err := h.courseStructureService.UploadExcel(openedFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	req := models.CourseStructureRequest{
		CourseID: courseID,
		Detail:   detail,
	}

	created, err := h.courseStructureService.CreateCourseStructure(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

func (h *CourseStructureHandler) UpdateCourseStructure(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course_structure_id"})
		return
	}

	courseID := c.PostForm("course_id")

	file, err := c.FormFile("detail")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer openedFile.Close()

	detail, err := h.courseStructureService.UploadExcel(openedFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	req := models.CourseStructureRequest{
		CourseID: courseID,
		Detail:   detail,
	}

	updated, err := h.courseStructureService.UpdateCourseStructure(id, req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "course_structure not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (h *CourseStructureHandler) DeleteCourseStructure(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course_structure ID"})
		return
	}

	err = h.courseStructureService.DeleteCourseStructure(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "course_structure not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "course_structure deleted successfully"})
}
