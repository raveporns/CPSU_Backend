package models

type Courses struct {
	CourseID      string `json:"course_id"`
	Degree        string `json:"degree"`
	Major         string `json:"major"`
	Year          int    `json:"year"`
	ThaiCourse    string `json:"thai_course"`
	EngCourse     string `json:"eng_course"`
	ThaiDegree    string `json:"thai_degree"`
	EngDegree     string `json:"eng_degree"`
	AdmissionReq  string `json:"admission_req"`
	GraduationReq string `json:"graduation_req"`
	Philosophy    string `json:"philosophy"`
	Objective     string `json:"objective"`
	Tuition       string `json:"tuition"`
	Credits       string `json:"credits"`
	CareerPathsID int    `json:"career_paths_id"`
	CareerPaths   string `json:"career_paths"`
	PloID         int    `json:"plo_id"`
	PLO           string `json:"plo"`
	DetailURL     string `json:"detail_url"`
	Status        string `json:"status"`
}

type CoursesQueryParam struct {
	Search string `form:"search"`
	Limit  int    `form:"limit"`
	Degree string `form:"degree"`
	Major  string `form:"major"`
	Year   int    `form:"year"`
	Sort   string `form:"sort"`
	Status string `form:"status"`
	Order  string `form:"order"`
}

type CoursesRequest struct {
	CourseID      string `json:"course_id"`
	Degree        string `json:"degree"`
	Major         string `json:"major"`
	Year          int    `json:"year"`
	ThaiCourse    string `json:"thai_course"`
	EngCourse     string `json:"eng_course"`
	ThaiDegree    string `json:"thai_degree"`
	EngDegree     string `json:"eng_degree"`
	AdmissionReq  string `json:"admission_req"`
	GraduationReq string `json:"graduation_req"`
	Philosophy    string `json:"philosophy"`
	Objective     string `json:"objective"`
	Tuition       string `json:"tuition"`
	Credits       string `json:"credits"`
	CareerPathsID int    `json:"career_paths_id"`
	CareerPaths   string `json:"career_paths"`
	PloID         int    `json:"plo_id"`
	PLO           string `json:"plo"`
	DetailURL     string `json:"detail_url"`
	Status        string `json:"status"`
}
