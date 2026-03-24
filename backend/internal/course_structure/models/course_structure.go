package models

type CourseStructure struct {
	CourseStructureID int    `json:"course_structure_id"`
	CourseID          string `json:"course_id"`
	ThaiCourse        string `json:"thai_course"`
	Detail            string `json:"detail"`
}

type CourseStructureQueryParam struct {
	Search   string `form:"search"`
	Limit    int    `form:"limit"`
	CourseID string `form:"course_id"`
	Sort     string `form:"sort"`
	Order    string `form:"order"`
}

type CourseStructureRequest struct {
	CourseID string `json:"course_id"`
	Detail   string `json:"detail"`
}
