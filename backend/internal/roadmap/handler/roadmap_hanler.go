package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"cpsu/internal/roadmap/models"
	"cpsu/internal/roadmap/service"

	"github.com/gin-gonic/gin"
)

type RoadmapHandler struct {
	roadmapService service.RoadmapService
}

func NewRoadmapHandler(roadmapService service.RoadmapService) *RoadmapHandler {
	return &RoadmapHandler{roadmapService: roadmapService}
}

func (h *RoadmapHandler) GetAllRoadmap(c *gin.Context) {
	var param models.RoadmapQueryParam
	if err := c.BindQuery(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roadmaps, err := h.roadmapService.GetAllRoadmap(param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, roadmaps)
}

func (h *RoadmapHandler) GetRoadmapByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid roadmap ID"})
		return
	}

	roadmap, err := h.roadmapService.GetRoadmapByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "roadmap not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, roadmap)
}

func (h *RoadmapHandler) CreateRoadmap(c *gin.Context) {
	courseID := c.PostForm("course_id")
	if courseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course ID"})
		return
	}

	file, err := c.FormFile("roadmap_url")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "roadmap_url is required"})
		return
	}

	created, err := h.roadmapService.CreateRoadmap(courseID, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

func (h *RoadmapHandler) DeleteRoadmap(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid roadmap ID"})
		return
	}

	err = h.roadmapService.DeleteRoadmap(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "roadmap not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "roadmap deleted successfully"})
}
