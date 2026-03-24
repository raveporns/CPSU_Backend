package models

type Admission struct {
	AdmissionID int    `json:"admission_id"`
	Round       string `json:"round"`
	Detail      string `json:"detail"`
	FileImage   string `json:"file_image"`
}

type AdmissionQueryParam struct {
	Search string `form:"search"`
	Limit  int    `form:"limit"`
	Round  string `form:"round"`
	Sort   string `form:"sort"`
	Order  string `form:"order"`
}

type AdmissionRequest struct {
	Round     string `json:"round"`
	Detail    string `json:"detail"`
	FileImage string `json:"file_image"`
}
