package service

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"

	"cpsu/internal/news/models"
	"cpsu/internal/news/repository"
	newsrepo "cpsu/internal/news/repository"

	authrepo "cpsu/internal/auth/repository"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type NewsService interface {
	GetAllNews(param models.NewsQueryParam) ([]models.News, error)
	GetNewsByID(id int) (*models.News, error)
	CreateNews(title, content string, typeID int, typeName, detailURL string, coverImage *multipart.FileHeader, images []*multipart.FileHeader, userID int, ip string, userAgent string) (*models.News, error)
	UpdateNews(id int, title, content string, type_id int, typeName, detailURL string, coverImage *multipart.FileHeader, images []*multipart.FileHeader, userID int, ip string, userAgent string) (*models.News, error)
	DeleteNews(id int, userID int, ip string, userAgent string) error
}

type newsService struct {
	repo        newsrepo.NewsRepository
	auditRepo   *authrepo.AuditRepository
	minioClient *minio.Client
	bucket      string
	publicBase  string
}

func NewNewsService(
	repo newsrepo.NewsRepository,
	auditRepo *authrepo.AuditRepository,
	endpoint string,
	accessKey string,
	secretKey string,
	bucket string,
	useSSL bool,
	publicBaseURL string,
) NewsService {

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		panic(err)
	}

	return &newsService{
		repo:        repo,
		auditRepo:   auditRepo,
		minioClient: client,
		bucket:      bucket,
		publicBase:  publicBaseURL,
	}
}

func (s *newsService) GetAllNews(param models.NewsQueryParam) ([]models.News, error) {
	if param.Sort == "" {
		param.Sort = "created_at"
	}
	if param.Order == "" {
		param.Order = "desc"
	}
	newsList, err := s.repo.GetAllNews(param)
	if err != nil {
		return nil, err
	}
	if len(newsList) == 0 && param.TypeID > 0 {
		return nil, errors.New("news type not found")
	}

	return newsList, nil
}

func (s *newsService) GetNewsByID(id int) (*models.News, error) {
	return s.repo.GetNewsByID(id)
}

func (s *newsService) CreateNews(title, content string, typeID int, typeName, detailURL string, coverImage *multipart.FileHeader, images []*multipart.FileHeader, userID int, ip string, userAgent string) (*models.News, error) {
	if strings.TrimSpace(title) == "" || strings.TrimSpace(content) == "" {
		return nil, errors.New("title and content are required")
	}

	var uploadedFlies []string
	for _, fileHeader := range images {
		url, err := s.UploadImages(fileHeader)
		if err != nil {
			return nil, err
		}
		uploadedFlies = append(uploadedFlies, url)
	}

	var coverURL string
	if coverImage != nil {
		url, err := s.UploadImages(coverImage)
		if err != nil {
			return nil, err
		}
		coverURL = url
	}

	newsReq := &models.NewsRequest{
		Title:      title,
		Content:    content,
		TypeID:     typeID,
		DetailURL:  detailURL,
		CoverImage: coverURL,
	}

	for _, url := range uploadedFlies {
		newsReq.Images = append(newsReq.Images, models.NewsImages{FileImage: url})
	}

	created, err := s.repo.CreateNews(newsReq)
	if err != nil {
		return nil, err
	}

	_ = s.auditRepo.LogAudit(
		userID, "create", "news", strconv.Itoa(created.NewsID),
		map[string]interface{}{
			"title": created.Title,
		},
		ip,
		userAgent,
	)

	return s.repo.GetNewsByID(created.NewsID)
}

func (s *newsService) UpdateNews(id int, title, content string, typeID int, typeName, detailURL string, coverImage *multipart.FileHeader, images []*multipart.FileHeader, userID int, ip string, userAgent string) (*models.News, error) {
	if strings.TrimSpace(title) == "" || strings.TrimSpace(content) == "" {
		return nil, errors.New("title and content are required")
	}

	var uploadedFlies []string
	for _, fileHeader := range images {
		url, err := s.UploadImages(fileHeader)
		if err != nil {
			return nil, err
		}
		uploadedFlies = append(uploadedFlies, url)
	}

	var coverURL string
	if coverImage != nil {
		url, err := s.UploadImages(coverImage)
		if err != nil {
			return nil, err
		}
		coverURL = url
	}

	newsReq := &models.NewsRequest{
		Title:      title,
		Content:    content,
		TypeID:     typeID,
		DetailURL:  detailURL,
		CoverImage: coverURL,
	}

	for _, url := range uploadedFlies {
		newsReq.Images = append(newsReq.Images, models.NewsImages{FileImage: url})
	}

	_, err := s.repo.UpdateNews(id, newsReq)
	if err != nil {
		return nil, err
	}

	_, err = s.repo.UpdateNewsImages(id, repository.ImagesAsStrings(newsReq.Images))
	if err != nil {
		return nil, err
	}

	_ = s.auditRepo.LogAudit(
		userID,
		"update",
		"news",
		strconv.Itoa(id),
		map[string]interface{}{
			"title": title,
		},
		ip,
		userAgent,
	)

	return s.repo.GetNewsByID(id)
}

func (s *newsService) DeleteNews(id int, userID int, ip string, userAgent string) error {

	err := s.repo.DeleteNews(id)
	if err != nil {
		return err
	}

	err = s.auditRepo.LogAudit(
		userID, "delete", "news", strconv.Itoa(id),
		map[string]interface{}{}, ip, userAgent,
	)

	return nil
}
func (s *newsService) UploadImages(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	ext := filepath.Ext(fileHeader.Filename)
	objectName := fmt.Sprintf(
		"news/%s%s",
		uuid.New().String(),
		ext,
	)

	_, err = s.minioClient.PutObject(
		context.Background(),
		s.bucket,
		objectName,
		file,
		fileHeader.Size,
		minio.PutObjectOptions{
			ContentType: fileHeader.Header.Get("Content-Type"),
		},
	)
	if err != nil {
		return "", err
	}

	imageURL := fmt.Sprintf(
		"%s/%s/%s",
		s.publicBase,
		s.bucket,
		objectName,
	)

	return imageURL, nil
}
