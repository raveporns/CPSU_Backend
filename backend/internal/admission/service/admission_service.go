package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strconv"

	"cpsu/internal/admission/models"
	"cpsu/internal/admission/repository"

	authrepo "cpsu/internal/auth/repository"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type AdmissionService interface {
	GetAllAdmission(param models.AdmissionQueryParam) ([]models.Admission, error)
	GetAdmissionByID(id int) (*models.Admission, error)
	CreateAdmission(req models.AdmissionRequest, fileImage *multipart.FileHeader, userID int, ip string, userAgent string) (*models.Admission, error)
	UpdateAdmission(id int, req models.AdmissionRequest, fileImage *multipart.FileHeader, userID int, ip string, userAgent string) (*models.Admission, error)
	DeleteAdmission(id int, userID int, ip string, userAgent string) error
}

type admissionService struct {
	repo        repository.AdmissionRepository
	auditRepo   *authrepo.AuditRepository
	minioClient *minio.Client
	bucket      string
	publicBase  string
}

func NewAdmissionService(
	repo repository.AdmissionRepository,
	auditRepo *authrepo.AuditRepository,
	endpoint string,
	accessKey string,
	secretKey string,
	bucket string,
	useSSL bool,
	publicBaseURL string,
) AdmissionService {

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		panic(err)
	}

	return &admissionService{
		repo:        repo,
		auditRepo:   auditRepo,
		minioClient: client,
		bucket:      bucket,
		publicBase:  publicBaseURL,
	}
}

func (s *admissionService) GetAllAdmission(param models.AdmissionQueryParam) ([]models.Admission, error) {
	return s.repo.GetAllAdmission(param)
}

func (s *admissionService) GetAdmissionByID(id int) (*models.Admission, error) {
	return s.repo.GetAdmissionByID(id)
}

func (s *admissionService) CreateAdmission(req models.AdmissionRequest, fileImage *multipart.FileHeader, userID int, ip string, userAgent string) (*models.Admission, error) {

	if fileImage != nil {
		url, err := s.uploadFile(fileImage)
		if err != nil {
			return nil, err
		}
		req.FileImage = url
	}

	created, err := s.repo.CreateAdmission(req)
	if err != nil {
		return nil, err
	}

	err = s.auditRepo.LogAudit(
		userID, "create", "admission",
		strconv.Itoa(created.AdmissionID),
		map[string]interface{}{
			"round": created.Round,
		},
		ip, userAgent,
	)

	return created, nil
}

func (s *admissionService) UpdateAdmission(id int, req models.AdmissionRequest, fileImage *multipart.FileHeader, userID int, ip string, userAgent string) (*models.Admission, error) {

	if fileImage != nil {
		url, err := s.uploadFile(fileImage)
		if err != nil {
			return nil, err
		}
		req.FileImage = url
	}

	updated, err := s.repo.UpdateAdmission(id, req)
	if err != nil {
		return nil, err
	}

	err = s.auditRepo.LogAudit(
		userID, "update", "admission", strconv.Itoa(id),
		map[string]interface{}{
			"round": req.Round,
		},
		ip, userAgent,
	)

	return updated, nil
}

func (s *admissionService) DeleteAdmission(id int, userID int, ip string, userAgent string) error {

	err := s.repo.DeleteAdmission(id)
	if err != nil {
		return err
	}

	err = s.auditRepo.LogAudit(
		userID, "delete", "admission", strconv.Itoa(id),
		map[string]interface{}{},
		ip, userAgent,
	)

	return nil
}
func (s *admissionService) uploadFile(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	ext := filepath.Ext(fileHeader.Filename)
	objectName := fmt.Sprintf(
		"admission/%s%s",
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
