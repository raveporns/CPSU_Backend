package repository

import (
	"database/sql"
	"strconv"
	"strings"

	"cpsu/internal/calendar/models"
)

type CalendarRepository interface {
	GetAllCalendars(param models.CalendarQueryParam) ([]models.Calendar, error)
	GetCalendarByID(id int) (*models.Calendar, error)
	CreateCalendar(req *models.CalendarRequest) (*models.Calendar, error)
	UpdateCalendar(req *models.CalendarRequest) (*models.Calendar, error)
	DeleteCalendar(id int) error
}

type calendarRepository struct {
	db *sql.DB
}

func NewCalendarRepository(db *sql.DB) CalendarRepository {
	return &calendarRepository{db: db}
}

func (r *calendarRepository) GetAllCalendars(param models.CalendarQueryParam) ([]models.Calendar, error) {
	query := `
		SELECT calendar_id, title, detail, start_date, end_date
		FROM calendar
	`

	sort := "calendar_id"
	if param.Sort != "" {
		sort = param.Sort
	}

	order := "ASC"
	if strings.ToUpper(param.Order) == "DESC" {
		order = "DESC"
	}

	query += " ORDER BY " + sort + " " + order

	if param.Limit > 0 {
		query += " LIMIT " + strconv.Itoa(param.Limit)
	}

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var calendars []models.Calendar
	for rows.Next() {
		var calendar models.Calendar
		err := rows.Scan(&calendar.CalenderID, &calendar.Title, &calendar.Detail, &calendar.StartDate, &calendar.EndDate)
		if err != nil {
			return nil, err
		}
		calendars = append(calendars, calendar)
	}

	return calendars, nil
}

func (r *calendarRepository) GetCalendarByID(id int) (*models.Calendar, error) {
	query := `
		SELECT calendar_id, title, detail, start_date, end_date
		FROM calendar
		WHERE calendar_id = $1
	`
	row := r.db.QueryRow(query, id)

	var calendar models.Calendar
	err := row.Scan(&calendar.CalenderID, &calendar.Title, &calendar.Detail, &calendar.StartDate, &calendar.EndDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return &calendar, nil
}

func (r *calendarRepository) CreateCalendar(req *models.CalendarRequest) (*models.Calendar, error) {
	query := `
		INSERT INTO calendar (title, detail, start_date, end_date)
		VALUES ($1, $2, $3, $4)
		RETURNING calendar_id
	`

	var calendar models.Calendar
	err := r.db.QueryRow(query, req.Title, req.Detail, req.StartDate, req.EndDate).Scan(&calendar.CalenderID)
	if err != nil {
		return nil, err
	}

	calendar.Title = req.Title
	calendar.Detail = req.Detail
	calendar.StartDate = req.StartDate
	calendar.EndDate = req.EndDate

	return &calendar, nil
}

func (r *calendarRepository) UpdateCalendar(req *models.CalendarRequest) (*models.Calendar, error) {
	query := `
		UPDATE calendar
		SET title = $1, detail = $2, start_date = $3, end_date = $4
		WHERE calendar_id = $5
		RETURNING calendar_id
	`

	var calendar models.Calendar
	err := r.db.QueryRow(query, req.Title, req.Detail, req.StartDate, req.EndDate, req.CalenderID).Scan(&calendar.CalenderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	calendar.Title = req.Title
	calendar.Detail = req.Detail
	calendar.StartDate = req.StartDate
	calendar.EndDate = req.EndDate

	return &calendar, nil
}

func (r *calendarRepository) DeleteCalendar(id int) error {
	result, err := r.db.Exec("DELETE FROM calendar WHERE calendar_id = $1", id)
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
