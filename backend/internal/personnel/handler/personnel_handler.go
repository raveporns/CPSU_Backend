package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"cpsu/internal/personnel/models"
	"cpsu/internal/personnel/service"

	"cpsu/internal/auth/repository"

	"github.com/gin-gonic/gin"
)

type PersonnelHandler struct {
	personnelService service.PersonnelService
	auditRepo        *repository.AuditRepository
}

func NewPersonnelHandler(personnelService service.PersonnelService) *PersonnelHandler {
	return &PersonnelHandler{personnelService: personnelService}
}

func (h *PersonnelHandler) GetAllPersonnels(c *gin.Context) {
	var param models.PersonnelQueryParam
	if err := c.BindQuery(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameter"})
		return
	}

	personnel, err := h.personnelService.GetAllPersonnels(param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, personnel)
}

func (h *PersonnelHandler) GetPersonnelByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	personnel, err := h.personnelService.GetPersonnelByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"personnel": personnel})
}

func (h *PersonnelHandler) CreatePersonnel(c *gin.Context) {

	req := models.PersonnelRequest{
		TypePersonnel:          c.PostForm("type_personnel"),
		DepartmentPositionID:   intPtr(c.PostForm("department_position_id")),
		DepartmentPositionName: strPtr(c.PostForm("department_position_name")),
		ThaiAcademicPosition:   strPtr(c.PostForm("thai_academic_position")),
		EngAcademicPosition:    strPtr(c.PostForm("eng_academic_position")),
		ThaiName:               c.PostForm("thai_name"),
		EngName:                c.PostForm("eng_name"),
		Education:              strPtr(c.PostForm("education")),
		RelatedFields:          strPtr(c.PostForm("related_fields")),
		Email:                  strPtr(c.PostForm("email")),
		Website:                strPtr(c.PostForm("website")),
		ScopusID:               strPtr(c.PostForm("scopus_id")),
		AcademicPositionID:     intPtr(c.PostForm("academic_position_id")),
	}
	fileImage, err := c.FormFile("file_image")
	if err != nil {
		fileImage = nil
	}

	userID := c.GetInt("user_id")
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	createdPersonnel, err := h.personnelService.CreatePersonnel(req, fileImage, userID, ip, userAgent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdPersonnel)
}

func (h *PersonnelHandler) UpdatePersonnel(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid personnel ID"})
		return
	}

	req := models.PersonnelRequest{
		TypePersonnel:          c.PostForm("type_personnel"),
		DepartmentPositionID:   intPtr(c.PostForm("department_position_id")),
		DepartmentPositionName: strPtr(c.PostForm("department_position_name")),
		ThaiAcademicPosition:   strPtr(c.PostForm("thai_academic_position")),
		EngAcademicPosition:    strPtr(c.PostForm("eng_academic_position")),
		ThaiName:               c.PostForm("thai_name"),
		EngName:                c.PostForm("eng_name"),
		Education:              strPtr(c.PostForm("education")),
		RelatedFields:          strPtr(c.PostForm("related_fields")),
		Email:                  strPtr(c.PostForm("email")),
		Website:                strPtr(c.PostForm("website")),
		ScopusID:               strPtr(c.PostForm("scopus_id")),
		AcademicPositionID:     intPtr(c.PostForm("academic_position_id")),
	}

	fileImage, err := c.FormFile("file_image")
	if err != nil {
		fileImage = nil
	}

	userID := c.GetInt("user_id")
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	updatedPersonnel, err := h.personnelService.UpdatePersonnel(id, req, fileImage, userID, ip, userAgent)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "personnel ID not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, updatedPersonnel)
}

func (h *PersonnelHandler) UpdateTeacher(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid personnel ID"})
		return
	}

	req := models.TeacherRequest{
		ThaiName:      c.PostForm("thai_name"),
		EngName:       c.PostForm("eng_name"),
		Education:     strPtr(c.PostForm("education")),
		RelatedFields: strPtr(c.PostForm("related_fields")),
		Email:         strPtr(c.PostForm("email")),
		Website:       strPtr(c.PostForm("website")),
		ScopusID:      strPtr(c.PostForm("scopus_id")),
	}

	fileImage, err := c.FormFile("file_image")
	if err != nil {
		fileImage = nil
	}

	userID := c.GetInt("user_id")
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	updatedTeacher, err := h.personnelService.UpdateTeacher(id, req, fileImage, userID, ip, userAgent)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "teacher ID not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, updatedTeacher)
}

func (h *PersonnelHandler) DeletePersonnel(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "personnel ID required"})
		return
	}

	userID := c.GetInt("user_id")
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	err = h.personnelService.DeletePersonnel(id, userID, ip, userAgent)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "personnel not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "personnel deleted successfully"})
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func intPtr(s string) *int {
	if s == "" {
		return nil
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}
	return &val
}

func (h *PersonnelHandler) GetResearchfromScopus(c *gin.Context) {
	pidStr := c.Query("personnel_id")
	if pidStr != "" {
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid personnel_id"})
			return
		}
		rs, err := h.personnelService.SyncResearch(pid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"personnel_id": pid, "researches": rs})
		return
	}

	count, err := h.personnelService.SyncAllFromScopus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data retrieved successfully", "processed_personnels": count})
}

func (h *PersonnelHandler) GetAllResearch(c *gin.Context) {
	var param models.ResearchQueryParam
	if err := c.ShouldBindQuery(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rs, err := h.personnelService.GetAllResearch(param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rs)
}
