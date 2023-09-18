package models

type EnrollmentCore struct {
	Created  string `json:"Created"`
	Mode     string `json:"Mode"`
	IsActive bool   `json:"IsActive"`
	User     string `json:"User"`
	CourseID string `json:"Course_ID"`
}
type EnrollmentsListCore struct {
	Next     string            `json:"Next"`
	Previous string            `json:"Previous"`
	Results  []*EnrollmentCore `json:"Results"`
}
