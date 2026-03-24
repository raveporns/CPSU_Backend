package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"cpsu/internal/admission/models"
	"cpsu/internal/admission/service"

	"cpsu/internal/auth/repository"

	"github.com/gin-gonic/gin"
)

type AdmissionHandler struct {
	admissionService service.AdmissionService
	auditRepo        *repository.AuditRepository
}

func NewAdmissionHandler(admissionService service.AdmissionService) *AdmissionHandler {
	return &AdmissionHandler{admissionService: admissionService}
}

func (h *AdmissionHandler) GetAllAdmission(c *gin.Context) {
	var param models.AdmissionQueryParam
	if err := c.ShouldBindQuery(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameter"})
		return
	}

	admissions, err := h.admissionService.GetAllAdmission(param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, admissions)
}

func (h *AdmissionHandler) GetAdmissionByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid admission ID"})
		return
	}

	admission, err := h.admissionService.GetAdmissionByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "admission not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, admission)
}

func (h *AdmissionHandler) CreateAdmission(c *gin.Context) {
	req := models.AdmissionRequest{
		Round:  c.PostForm("round"),
		Detail: c.PostForm("detail"),
	}

	fileImage, _ := c.FormFile("file_image")

	userID := c.GetInt("user_id")
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	created, err := h.admissionService.CreateAdmission(req, fileImage, userID, ip, userAgent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

func (h *AdmissionHandler) UpdateAdmission(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid admission ID"})
		return
	}

	req := models.AdmissionRequest{
		Round:  c.PostForm("round"),
		Detail: c.PostForm("detail"),
	}

	fileImage, _ := c.FormFile("file_image")

	userID := c.GetInt("user_id")
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	updated, err := h.admissionService.UpdateAdmission(id, req, fileImage, userID, ip, userAgent)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "admission not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (h *AdmissionHandler) DeleteAdmission(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "admission ID required"})
		return
	}

	userID := c.GetInt("user_id")
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	err = h.admissionService.DeleteAdmission(id, userID, ip, userAgent)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "admission not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "admission deleted successfully"})
}
