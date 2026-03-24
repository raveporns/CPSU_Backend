package handler

import (
	"database/sql"
	"errors"
	"mime/multipart"
	"net/http"
	"strconv"

	"cpsu/internal/news/models"
	"cpsu/internal/news/service"

	"cpsu/internal/auth/repository"

	"github.com/gin-gonic/gin"
)

type NewsHandler struct {
	newsService service.NewsService
	auditRepo   *repository.AuditRepository
}

func NewNewsHandler(newsService service.NewsService) *NewsHandler {
	return &NewsHandler{newsService: newsService}
}

func (h *NewsHandler) GetAllNews(c *gin.Context) {
	var param models.NewsQueryParam
	if err := c.BindQuery(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameter"})
		return
	}

	newsList, err := h.newsService.GetAllNews(param)
	if err != nil {
		if err.Error() == "news type not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get news"})
		}
		return
	}
	c.JSON(http.StatusOK, newsList)
}

func (h *NewsHandler) GetNewsByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid news ID"})
		return
	}

	news, err := h.newsService.GetNewsByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "news ID not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, news)
}

func (h *NewsHandler) CreateNews(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")
	typeIDStr := c.PostForm("type_id")
	detailURL := c.PostForm("detail_url")

	coverImage, _ := c.FormFile("cover_image")

	typeID, err := strconv.Atoi(typeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid type ID"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form"})
		return
	}

	var fileImages []*multipart.FileHeader
	if form != nil {
		fileImages = form.File["images"]
	}

	userID := c.GetInt("user_id")
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	created, err := h.newsService.CreateNews(
		title, content, typeID, "", detailURL, coverImage, fileImages, userID, ip, userAgent,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

func (h *NewsHandler) UpdateNews(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid news ID"})
		return
	}

	title := c.PostForm("title")
	content := c.PostForm("content")
	typeIDStr := c.PostForm("type_id")
	detailURL := c.PostForm("detail_url")
	coverImage, _ := c.FormFile("cover_image")

	typeID, err := strconv.Atoi(typeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid type ID"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form"})
		return
	}

	var fileImages []*multipart.FileHeader
	if form != nil {
		fileImages = form.File["images"]
	}

	userID := c.GetInt("user_id")
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	updated, err := h.newsService.UpdateNews(id, title, content, typeID, "", detailURL, coverImage, fileImages, userID, ip, userAgent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (h *NewsHandler) DeleteNews(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid news ID"})
		return
	}

	userID := c.GetInt("user_id")
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	err = h.newsService.DeleteNews(id, userID, ip, userAgent)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "news ID not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "News deleted successfully"})
}
