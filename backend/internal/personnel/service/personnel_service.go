package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	authrepo "cpsu/internal/auth/repository"
	"cpsu/internal/personnel/models"
	"cpsu/internal/personnel/repository"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type PersonnelService interface {
	GetAllPersonnels(param models.PersonnelQueryParam) ([]models.Personnels, error)
	GetPersonnelByID(id int) (*models.Personnels, error)
	CreatePersonnel(req models.PersonnelRequest, fileImage *multipart.FileHeader, userID int, ip string, userAgent string) (*models.Personnels, error)
	UpdatePersonnel(id int, req models.PersonnelRequest, fileImage *multipart.FileHeader, userID int, ip string, userAgent string) (*models.Personnels, error)
	UpdateTeacher(id int, req models.TeacherRequest, fileImage *multipart.FileHeader, userID int, ip string, userAgent string) (*models.Personnels, error)
	DeletePersonnel(id int, userID int, ip string, userAgent string) error
	SyncResearch(personnelID int) ([]models.Research, error)
	GetResearchFromScopus(scopusID string) ([]models.Research, error)
	SyncAllFromScopus() (int, error)
	GetAllResearch(param models.ResearchQueryParam) ([]models.Research, error)
}

type personnelService struct {
	repo        repository.PersonnelRepository
	auditRepo   *authrepo.AuditRepository
	minioClient *minio.Client
	bucket      string
	publicBase  string
}

func NewPersonnelService(
	repo repository.PersonnelRepository,
	auditRepo *authrepo.AuditRepository,
	endpoint string,
	accessKey string,
	secretKey string,
	bucket string,
	useSSL bool,
	publicBaseURL string,
) PersonnelService {

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		panic(err)
	}

	return &personnelService{
		repo:        repo,
		auditRepo:   auditRepo,
		minioClient: client,
		bucket:      bucket,
		publicBase:  publicBaseURL,
	}
}

func (s *personnelService) GetAllPersonnels(param models.PersonnelQueryParam) ([]models.Personnels, error) {
	return s.repo.GetAllPersonnels(param)
}

func (s *personnelService) GetPersonnelByID(id int) (*models.Personnels, error) {
	return s.repo.GetPersonnelByID(id)
}

func (s *personnelService) CreatePersonnel(req models.PersonnelRequest, fileImage *multipart.FileHeader, userID int, ip string, userAgent string) (*models.Personnels, error) {

	if fileImage != nil {
		url, err := s.uploadFile(fileImage)
		if err != nil {
			return nil, err
		}
		req.FileImage = url
	}

	created, err := s.repo.CreatePersonnel(req)
	if err != nil {
		return nil, err
	}

	err = s.auditRepo.LogAudit(
		userID, "create", "personnel",
		strconv.Itoa(created.PersonnelID),
		map[string]interface{}{
			"thai_name": created.ThaiName,
		},
		ip, userAgent,
	)

	return created, nil
}

func (s *personnelService) UpdatePersonnel(id int, req models.PersonnelRequest, fileImage *multipart.FileHeader, userID int, ip string, userAgent string) (*models.Personnels, error) {

	if fileImage != nil {
		url, err := s.uploadFile(fileImage)
		if err != nil {
			return nil, err
		}
		req.FileImage = url
	}

	updated, err := s.repo.UpdatePersonnel(id, req)
	if err != nil {
		return nil, err
	}

	err = s.auditRepo.LogAudit(
		userID, "update", "personnel", strconv.Itoa(id),
		map[string]interface{}{
			"thai_name": req.ThaiName,
		},
		ip, userAgent,
	)

	return updated, nil
}

func (s *personnelService) UpdateTeacher(id int, req models.TeacherRequest, fileImage *multipart.FileHeader, userID int, ip string, userAgent string) (*models.Personnels, error) {

	if fileImage != nil {
		url, err := s.uploadFile(fileImage)
		if err != nil {
			return nil, err
		}
		req.FileImage = url
	}

	existing, err := s.repo.GetPersonnelByID(id)
	if err != nil {
		return nil, err
	}

	personnelReq := models.PersonnelRequest{
		TypePersonnel:        existing.TypePersonnel,
		DepartmentPositionID: &existing.DepartmentPositionID,
		AcademicPositionID:   existing.AcademicPositionID,
		ThaiAcademicPosition: existing.ThaiAcademicPosition,
		EngAcademicPosition:  existing.EngAcademicPosition,
		ThaiName:             req.ThaiName,
		EngName:              req.EngName,
		Education:            req.Education,
		RelatedFields:        req.RelatedFields,
		Email:                req.Email,
		Website:              req.Website,
		FileImage:            req.FileImage,
		ScopusID:             req.ScopusID,
	}

	updated, err := s.repo.UpdatePersonnel(id, personnelReq)
	if err != nil {
		return nil, err
	}

	err = s.auditRepo.LogAudit(
		userID, "update", "personnel", strconv.Itoa(id),
		map[string]interface{}{
			"thai_name": req.ThaiName,
		},
		ip, userAgent,
	)

	return updated, nil
}

func (s *personnelService) DeletePersonnel(id int, userID int, ip string, userAgent string) error {

	err := s.repo.DeletePersonnel(id)
	if err != nil {
		return err
	}

	err = s.auditRepo.LogAudit(
		userID,
		"delete",
		"personnel",
		strconv.Itoa(id),
		map[string]interface{}{},
		ip,
		userAgent,
	)

	return nil
}

func (s *personnelService) uploadFile(fileHeader *multipart.FileHeader) (string, error) {
	if fileHeader == nil {
		return "", errors.New("file is nil")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	ext := filepath.Ext(fileHeader.Filename)
	objectName := fmt.Sprintf(
		"personnel/%s%s",
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

func (s *personnelService) SyncResearch(personnelID int) ([]models.Research, error) {
	scopusIDptr, err := s.repo.GetScopusIDByPersonnelID(personnelID)
	if err != nil {
		return nil, err
	}
	if scopusIDptr == nil || *scopusIDptr == "" {
		return nil, fmt.Errorf("personnel %d has no scopus_id", personnelID)
	}

	rs, err := s.GetResearchFromScopus(*scopusIDptr)
	if err != nil {
		return nil, err
	}
	for i := range rs {
		rs[i].PersonnelID = personnelID
		if rs[i].CreatedAt.IsZero() {
			rs[i].CreatedAt = time.Now()
		}
	}
	if err := s.repo.SaveResearch(personnelID, rs); err != nil {
		return nil, err
	}

	param := models.ResearchQueryParam{PersonnelID: personnelID}
	return s.repo.GetAllResearch(param)
}

func (s *personnelService) getAuthorsFromAbstract(eid string) []string {
	apiKey := os.Getenv("SCOPUS_API_KEY")

	url := fmt.Sprintf(
		"https://api.elsevier.com/content/abstract/scopus_id/%s",
		eid,
	)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-ELS-APIKey", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return []string{}
	}

	root, ok := data["abstracts-retrieval-response"].(map[string]interface{})
	if !ok {
		return []string{}
	}

	item, ok := root["item"].(map[string]interface{})
	if ok {

		bib := item["bibrecord"].(map[string]interface{})
		head := bib["head"].(map[string]interface{})

		ag := head["author-group"]

		var groups []interface{}

		if arr, ok := ag.([]interface{}); ok {
			groups = arr
		} else if obj, ok := ag.(map[string]interface{}); ok {
			groups = []interface{}{obj}
		}

		var authors []string

		for _, g := range groups {

			group := g.(map[string]interface{})
			a := group["author"]

			var authorList []interface{}

			if arr, ok := a.([]interface{}); ok {
				authorList = arr
			} else if obj, ok := a.(map[string]interface{}); ok {
				authorList = []interface{}{obj}
			}

			for _, x := range authorList {

				author := x.(map[string]interface{})
				pn := author["preferred-name"].(map[string]interface{})

				if name, ok := pn["ce:indexed-name"].(string); ok {
					authors = append(authors, name)
				} else {
					given := fmt.Sprint(pn["ce:given-name"])
					sur := fmt.Sprint(pn["ce:surname"])
					authors = append(authors, strings.TrimSpace(given+" "+sur))
				}
			}
		}

		if len(authors) > 0 {
			return authors
		}
	}

	if core, ok := root["coredata"].(map[string]interface{}); ok {
		if creator, ok := core["dc:creator"].(string); ok {
			return []string{creator}
		}
	}

	return []string{}
}

func (s *personnelService) GetResearchFromScopus(scopusID string) ([]models.Research, error) {
	apiKey := os.Getenv("SCOPUS_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("missing SCOPUS_API_KEY")
	}

	searchURL := fmt.Sprintf(
		"https://api.elsevier.com/content/search/scopus?query=AU-ID(%s)&apiKey=%s",
		scopusID, apiKey,
	)

	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("scopus search status %d", resp.StatusCode)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	sr, ok := data["search-results"].(map[string]interface{})
	if !ok {
		return []models.Research{}, nil
	}

	entries, ok := sr["entry"].([]interface{})
	if !ok || len(entries) == 0 {
		return []models.Research{}, nil
	}

	toPtr := func(v interface{}) *string {
		if v == nil {
			return nil
		}
		s := strings.TrimSpace(fmt.Sprint(v))
		if s == "" || s == "<nil>" {
			return nil
		}
		return &s
	}

	var researches []models.Research

	for _, e := range entries {
		item, ok := e.(map[string]interface{})
		if !ok {
			continue
		}

		year := 0
		if ds, ok := item["prism:coverDate"].(string); ok && len(ds) >= 4 {
			year, _ = strconv.Atoi(ds[:4])
		}

		cited := 0
		switch v := item["citedby-count"].(type) {
		case string:
			cited, _ = strconv.Atoi(v)
		case float64:
			cited = int(v)
		}

		rawID := fmt.Sprint(item["dc:identifier"])
		eid := strings.TrimPrefix(rawID, "SCOPUS_ID:")
		if eid == "" || eid == rawID {
			continue
		}

		fmt.Println("EID:", eid)

		authors := s.getAuthorsFromAbstract(eid)

		fmt.Println("Authors:", authors)

		researches = append(researches, models.Research{
			Title:     fmt.Sprint(item["dc:title"]),
			Journal:   fmt.Sprint(item["prism:publicationName"]),
			Year:      year,
			Volume:    toPtr(item["prism:volume"]),
			Issue:     toPtr(item["prism:issueIdentifier"]),
			Pages:     toPtr(item["prism:pageRange"]),
			DOI:       toPtr(item["prism:doi"]),
			Cited:     cited,
			Authors:   authors,
			CreatedAt: time.Now(),
		})
	}

	return researches, nil
}

func (s *personnelService) SyncAllFromScopus() (int, error) {
	personnels, err := s.repo.GetAllPersonnels(models.PersonnelQueryParam{})
	if err != nil {
		return 0, err
	}

	processed := 0

	for _, p := range personnels {
		if p.ScopusID == nil || *p.ScopusID == "" {
			continue
		}

		rs, err := s.GetResearchFromScopus(*p.ScopusID)
		if err != nil || len(rs) == 0 {
			continue
		}

		for i := range rs {
			rs[i].PersonnelID = p.PersonnelID
			if rs[i].CreatedAt.IsZero() {
				rs[i].CreatedAt = time.Now()
			}
		}

		if err := s.repo.SaveResearch(p.PersonnelID, rs); err == nil {
			processed++
		}
	}

	return processed, nil
}

func (s *personnelService) GetAllResearch(param models.ResearchQueryParam) ([]models.Research, error) {
	return s.repo.GetAllResearch(param)
}
