package service

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"

	"cpsu/internal/roadmap/models"
	"cpsu/internal/roadmap/repository"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type RoadmapService interface {
	GetAllRoadmap(param models.RoadmapQueryParam) ([]models.Roadmap, error)
	GetRoadmapByID(id int) (*models.Roadmap, error)
	CreateRoadmap(courseID string, file *multipart.FileHeader) (*models.Roadmap, error)
	DeleteRoadmap(id int) error
}

type roadmapService struct {
	repo        repository.RoadmapRepository
	minioClient *minio.Client
	bucket      string
	publicBase  string
}

func NewRoadmapService(
	repo repository.RoadmapRepository,
	endpoint string,
	accessKey string,
	secretKey string,
	bucket string,
	useSSL bool,
	publicBaseURL string,
) RoadmapService {

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		panic(err)
	}

	return &roadmapService{
		repo:        repo,
		minioClient: client,
		bucket:      bucket,
		publicBase:  publicBaseURL,
	}
}

func (s *roadmapService) GetAllRoadmap(param models.RoadmapQueryParam) ([]models.Roadmap, error) {
	return s.repo.GetAllRoadmap(param)
}

func (s *roadmapService) GetRoadmapByID(id int) (*models.Roadmap, error) {
	return s.repo.GetRoadmapByID(id)
}

func (s *roadmapService) CreateRoadmap(courseID string, file *multipart.FileHeader) (*models.Roadmap, error) {
	if courseID == "" {
		return nil, errors.New("course_id is required")
	}
	if file == nil {
		return nil, errors.New("roadmap image is required")
	}

	url, err := s.uploadFile(file)
	if err != nil {
		return nil, err
	}

	req := &models.RoadmapRequest{
		CourseID:   courseID,
		RoadmapURL: url,
	}

	return s.repo.CreateRoadmap(req)
}

func (s *roadmapService) DeleteRoadmap(id int) error {
	return s.repo.DeleteRoadmap(id)
}

func (s *roadmapService) uploadFile(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	ext := filepath.Ext(fileHeader.Filename)
	objectName := fmt.Sprintf(
		"roadmap/%s%s",
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
