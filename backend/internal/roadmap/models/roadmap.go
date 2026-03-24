package models

type Roadmap struct {
	RoadmapID  int    `json:"roadmap_id"`
	CourseID   string `json:"course_id"`
	ThaiCourse string `json:"thai_course"`
	RoadmapURL string `json:"roadmap_url"`
}

type RoadmapQueryParam struct {
	Search   string `form:"search"`
	Limit    int    `form:"limit"`
	CourseID string `form:"course_id"`
	Sort     string `form:"sort"`
	Order    string `form:"order"`
}

type RoadmapRequest struct {
	CourseID   string `json:"course_id"`
	RoadmapURL string `json:"roadmap_url"`
}
