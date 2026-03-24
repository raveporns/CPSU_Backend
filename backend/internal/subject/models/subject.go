package models

type Subjects struct {
	ID                int     `json:"id"`
	SubjectID         string  `json:"subject_id"`
	CourseID          string  `json:"course_id"`
	ThaiCourse        string  `json:"thai_course"`
	PlanType          string  `json:"plan_type"`
	Semester          string  `json:"semester"`
	ThaiSubject       string  `json:"thai_subject"`
	EngSubject        *string `json:"eng_subject,omitempty"`
	Credits           string  `json:"credits"`
	CompulsorySubject *string `json:"compulsory_subject,omitempty"`
	Condition         *string `json:"condition,omitempty"`
	DescriptionID     *string `json:"description_id,omitempty"`
	DescriptionThai   *string `json:"description_thai,omitempty"`
	DescriptionEng    *string `json:"description_eng,omitempty"`
	CloID             *string `json:"clo_id,omitempty"`
	CLO               *string `json:"clo,omitempty"`
}

type SubjectsQueryParam struct {
	Search    string `form:"search"`
	Limit     int    `form:"limit"`
	SubjectID string `form:"subject_id"`
	CourseID  string `form:"course_id"`
	PlanType  string `form:"plan_type"`
	Semester  string `form:"semester"`
	Sort      string `form:"sort"`
	Order     string `form:"order"`
}

type SubjectsRequest struct {
	SubjectID         string  `json:"subject_id"`
	CourseID          string  `json:"course_id"`
	PlanType          string  `json:"plan_type"`
	Semester          string  `json:"semester"`
	ThaiSubject       string  `json:"thai_subject"`
	EngSubject        *string `json:"eng_subject,omitempty"`
	Credits           string  `json:"credits"`
	CompulsorySubject *string `json:"compulsory_subject,omitempty"`
	Condition         *string `json:"condition,omitempty"`
	DescriptionID     *string `json:"description_id,omitempty"`
	DescriptionThai   *string `json:"description_thai,omitempty"`
	DescriptionEng    *string `json:"description_eng,omitempty"`
	CloID             *string `json:"clo_id,omitempty"`
	CLO               *string `json:"clo,omitempty"`
}
