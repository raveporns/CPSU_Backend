package models

import "time"

type Calendar struct {
	CalenderID int       `json:"calendar_id"`
	Title      string    `json:"title"`
	Detail     string    `json:"detail"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
}

type CalendarQueryParam struct {
	Search string `form:"search"`
	Limit  int    `form:"limit"`
	Sort   string `form:"sort"`
	Order  string `form:"order"`
}

type CalendarRequest struct {
	CalenderID int       `json:"calendar_id"`
	Title      string    `json:"title"`
	Detail     string    `json:"detail"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
}
