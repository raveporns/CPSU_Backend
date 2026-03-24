package models

import "time"

type News struct {
	NewsID     int          `json:"news_id"`
	Title      string       `json:"title"`
	Content    string       `json:"content"`
	TypeID     int          `json:"type_id"`
	TypeName   string       `json:"type_name"`
	DetailURL  string       `json:"detail_url"`
	CoverImage string       `json:"cover_image"`
	Images     []NewsImages `json:"images"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"update_at"`
}

type NewsImages struct {
	ImageID   int    `json:"image_id"`
	NewsID    int    `json:"news_id"`
	FileImage string `json:"file_image"`
}

type NewsQueryParam struct {
	Search string `form:"search"`
	Limit  int    `form:"limit"`
	TypeID int    `form:"type_id"`
	Sort   string `form:"sort"`
	Order  string `form:"order"`
}

type NewsRequest struct {
	Title      string       `json:"title"`
	Content    string       `json:"content"`
	TypeID     int          `json:"type_id"`
	DetailURL  string       `json:"detail_url"`
	CoverImage string       `json:"cover_image"`
	Images     []NewsImages `json:"images"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"update_at"`
}
