package repository

import (
	"database/sql"
	"strconv"
	"strings"

	"cpsu/internal/news/models"
)

type NewsRepository interface {
	GetAllNews(param models.NewsQueryParam) ([]models.News, error)
	GetNewsByID(id int) (*models.News, error)
	CreateNews(news *models.NewsRequest) (*models.News, error)
	UpdateNews(id int, news *models.NewsRequest) (*models.News, error)
	DeleteNews(id int) error
	AddNewsImages(newsID int, images []string) error
	UpdateNewsImages(newsID int, images []string) ([]models.NewsImages, error)
}

type newsRepository struct {
	db *sql.DB
}

func NewNewsRepository(db *sql.DB) NewsRepository {
	return &newsRepository{db: db}
}

func (r *newsRepository) GetAllNews(param models.NewsQueryParam) ([]models.News, error) {
	query := `
		SELECT n.news_id, n.title, n.content, nt.type_id, nt.type_name,
			n.detail_url, n.cover_image, ni.image_id, ni.file_image,
			n.created_at, n.updated_at
		FROM news n
		LEFT JOIN news_images ni ON n.news_id = ni.news_id
		LEFT JOIN news_types nt ON n.type_id = nt.type_id
	`

	conditions := []string{}
	args := []interface{}{}
	argIndex := 1

	if param.TypeID > 0 {
		conditions = append(conditions, "nt.type_id = $"+strconv.Itoa(argIndex))
		args = append(args, param.TypeID)
		argIndex++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	sort := "n.updated_at"
	if param.Sort != "" {
		sort = "n." + param.Sort
	}

	order := "ASC"
	if strings.ToUpper(param.Order) == "DESC" {
		order = "DESC"
	}

	query += " ORDER BY " + sort + " " + order

	if param.Limit > 0 {
		query += " LIMIT " + strconv.Itoa(param.Limit)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	newsMap := make(map[int]*models.News)
	var allNewsPt []*models.News

	for rows.Next() {
		var (
			newsID                                          int
			title, content, typeName, detailURL, coverImage string
			typeID, imageID                                 sql.NullInt64
			fileImage                                       sql.NullString
			createdAt, updatedAt                            sql.NullTime
		)

		err := rows.Scan(
			&newsID, &title, &content, &typeID, &typeName,
			&detailURL, &coverImage, &imageID, &fileImage,
			&createdAt, &updatedAt,
		)
		if err != nil {
			return nil, err
		}

		n, exists := newsMap[newsID]
		if !exists {
			n = &models.News{
				NewsID:     newsID,
				Title:      title,
				Content:    content,
				TypeID:     int(typeID.Int64),
				TypeName:   typeName,
				DetailURL:  detailURL,
				CoverImage: coverImage,
				CreatedAt:  createdAt.Time,
				UpdatedAt:  updatedAt.Time,
				Images:     []models.NewsImages{},
			}
			newsMap[newsID] = n
			allNewsPt = append(allNewsPt, n)
		}

		if imageID.Valid && fileImage.Valid {
			n.Images = append(n.Images, models.NewsImages{
				ImageID:   int(imageID.Int64),
				NewsID:    newsID,
				FileImage: fileImage.String,
			})
		}
	}

	var allNews []models.News
	for _, p := range allNewsPt {
		allNews = append(allNews, *p)
	}

	return allNews, nil
}

func (r *newsRepository) GetNewsByID(id int) (*models.News, error) {
	query := `
		SELECT n.news_id, n.title, n.content, nt.type_id, nt.type_name,
			n.detail_url, n.cover_image, ni.image_id, ni.file_image,
			n.created_at, n.updated_at
		FROM news n
		LEFT JOIN news_images ni ON n.news_id = ni.news_id
		LEFT JOIN news_types nt ON n.type_id = nt.type_id
		WHERE n.news_id = $1
	`
	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var news *models.News

	for rows.Next() {
		var (
			newsID                                          int
			title, content, typeName, detailURL, coverImage string
			typeID, imageID                                 sql.NullInt64
			fileImage                                       sql.NullString
			createdAt, updatedAt                            sql.NullTime
		)

		err := rows.Scan(
			&newsID, &title, &content, &typeID, &typeName,
			&detailURL, &coverImage, &imageID, &fileImage,
			&createdAt, &updatedAt,
		)
		if err != nil {
			return nil, err
		}

		if news == nil {
			news = &models.News{
				NewsID:     newsID,
				Title:      title,
				Content:    content,
				TypeID:     int(typeID.Int64),
				TypeName:   typeName,
				DetailURL:  detailURL,
				CoverImage: coverImage,
				CreatedAt:  createdAt.Time,
				UpdatedAt:  updatedAt.Time,
				Images:     []models.NewsImages{},
			}
		}

		if imageID.Valid && fileImage.Valid {
			news.Images = append(news.Images, models.NewsImages{
				ImageID:   int(imageID.Int64),
				NewsID:    newsID,
				FileImage: fileImage.String,
			})
		}
	}

	if news == nil {
		return nil, sql.ErrNoRows
	}

	return news, nil
}

func (r *newsRepository) CreateNews(req *models.NewsRequest) (*models.News, error) {
	query := `
		INSERT INTO news (title, content, type_id, detail_url, cover_image)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING news_id, created_at, updated_at
	`
	var news models.News
	err := r.db.QueryRow(
		query,
		req.Title, req.Content, req.TypeID, req.DetailURL, req.CoverImage,
	).Scan(&news.NewsID, &news.CreatedAt, &news.UpdatedAt)
	if err != nil {
		return nil, err
	}

	news.Title = req.Title
	news.Content = req.Content
	news.TypeID = req.TypeID
	news.DetailURL = req.DetailURL
	news.CoverImage = req.CoverImage
	news.Images = req.Images

	if len(req.Images) > 0 {
		err := r.AddNewsImages(news.NewsID, ImagesAsStrings(req.Images))
		if err != nil {
			return nil, err
		}
	}

	return &news, nil
}

func (r *newsRepository) UpdateNews(id int, newsreq *models.NewsRequest) (*models.News, error) {
	query := `
		UPDATE news
		SET title = $1, content = $2, type_id = $3, detail_url = $4, cover_image = $5, updated_at = NOW()
		WHERE news_id = $6
		RETURNING news_id, created_at, updated_at
	`
	var news models.News
	err := r.db.QueryRow(
		query,
		newsreq.Title, newsreq.Content, newsreq.TypeID, newsreq.DetailURL, newsreq.CoverImage, id,
	).Scan(&news.NewsID, &news.CreatedAt, &news.UpdatedAt)
	if err != nil {
		return nil, err
	}

	news.Title = newsreq.Title
	news.Content = newsreq.Content
	news.TypeID = newsreq.TypeID
	news.DetailURL = newsreq.DetailURL
	news.CoverImage = newsreq.CoverImage
	news.Images = newsreq.Images

	return &news, nil
}

func (r *newsRepository) DeleteNews(id int) error {
	result, err := r.db.Exec("DELETE FROM news WHERE news_id = $1", id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *newsRepository) AddNewsImages(newsID int, images []string) error {
	for _, img := range images {
		_, err := r.db.Exec("INSERT INTO news_images (news_id, file_image) VALUES ($1, $2)", newsID, img)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *newsRepository) UpdateNewsImages(newsID int, images []string) ([]models.NewsImages, error) {
	_, err := r.db.Exec("DELETE FROM news_images WHERE news_id = $1", newsID)
	if err != nil {
		return nil, err
	}

	for _, img := range images {
		_, err := r.db.Exec("INSERT INTO news_images (news_id, file_image) VALUES ($1, $2)", newsID, img)
		if err != nil {
			return nil, err
		}
	}

	rows, err := r.db.Query("SELECT image_id, news_id, file_image FROM news_images WHERE news_id = $1", newsID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.NewsImages
	for rows.Next() {
		var img models.NewsImages
		if err := rows.Scan(&img.ImageID, &img.NewsID, &img.FileImage); err != nil {
			return nil, err
		}
		result = append(result, img)
	}

	return result, nil
}

func (r *newsRepository) GetTypeNameByID(typeID int) (string, error) {
	var typeName string
	err := r.db.QueryRow("SELECT type_name FROM news_types WHERE type_id = $1", typeID).Scan(&typeName)
	if err != nil {
		return "", err
	}
	return typeName, nil
}

func ImagesAsStrings(images []models.NewsImages) []string {
	var files []string
	for _, img := range images {
		files = append(files, img.FileImage)
	}
	return files
}
