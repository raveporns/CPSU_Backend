package models

import "time"

type Personnels struct {
	PersonnelID            int        `json:"personnel_id"`
	TypePersonnel          string     `json:"type_personnel"`
	DepartmentPositionID   int        `json:"department_position_id"`
	DepartmentPositionName string     `json:"department_position_name"`
	AcademicPositionID     *int       `json:"academic_position_id,omitempty"`
	ThaiAcademicPosition   *string    `json:"thai_academic_position,omitempty"`
	EngAcademicPosition    *string    `json:"eng_academic_position,omitempty"`
	ThaiName               string     `json:"thai_name"`
	EngName                string     `json:"eng_name"`
	Education              *string    `json:"education,omitempty"`
	RelatedFields          *string    `json:"related_fields,omitempty"`
	Email                  *string    `json:"email,omitempty"`
	Website                *string    `json:"website,omitempty"`
	FileImage              string     `json:"file_image"`
	ScopusID               *string    `json:"scopus_id,omitempty"`
	Researches             []Research `json:"researches,omitempty"`
}

type PersonnelQueryParam struct {
	Search               string `form:"search"`
	Limit                int    `form:"limit"`
	TypePersonnel        string `form:"type_personnel"`
	DepartmentPositionID int    `form:"department_position_id"`
	AcademicPositionID   *int   `form:"academic_position_id"`
	Sort                 string `form:"sort"`
	Order                string `form:"order"`
}

type PersonnelRequest struct {
	TypePersonnel          string  `json:"type_personnel"`
	DepartmentPositionID   *int    `json:"department_position_id"`
	DepartmentPositionName *string `json:"department_position_name"`
	AcademicPositionID     *int    `json:"academic_position_id"`
	ThaiAcademicPosition   *string `json:"thai_academic_position"`
	EngAcademicPosition    *string `json:"eng_academic_position"`
	ThaiName               string  `json:"thai_name"`
	EngName                string  `json:"eng_name"`
	Education              *string `json:"education"`
	RelatedFields          *string `json:"related_fields"`
	Email                  *string `json:"email"`
	Website                *string `json:"website"`
	FileImage              string  `json:"file_image"`
	ScopusID               *string `json:"scopus_id"`
}

type TeacherRequest struct {
	ThaiName      string  `json:"thai_name"`
	EngName       string  `json:"eng_name"`
	Education     *string `json:"education"`
	RelatedFields *string `json:"related_fields"`
	Email         *string `json:"email"`
	Website       *string `json:"website"`
	FileImage     string  `json:"file_image"`
	ScopusID      *string `json:"scopus_id"`
}

type Research struct {
	ResearchID  int       `json:"research_id"`
	PersonnelID int       `json:"personnel_id"`
	ThaiName    string    `json:"thai_name,omitempty"`
	Authors     []string  `json:"authors"`
	Title       string    `json:"title"`
	Journal     string    `json:"journal"`
	Year        int       `json:"year"`
	Volume      *string   `json:"volume,omitempty"`
	Issue       *string   `json:"issue,omitempty"`
	Pages       *string   `json:"pages,omitempty"`
	DOI         *string   `json:"doi,omitempty"`
	Cited       int       `json:"cited"`
	CreatedAt   time.Time `json:"created_at"`
}

type ResearchAuthor struct {
	AuthorID    int    `json:"author_id"`
	ResearchID  int    `json:"research_id"`
	AuthorName  string `json:"author_name"`
	AuthorOrder int    `json:"author_order"`
}

type ResearchQueryParam struct {
	Search      string `form:"search"`
	Limit       int    `form:"limit"`
	PersonnelID int    `form:"personnel_id"`
	Sort        string `form:"sort"`
	Order       string `form:"order"`
}
