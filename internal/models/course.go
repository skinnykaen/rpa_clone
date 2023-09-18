package models

import (
	"time"
)

type PostEnrollmentHTTP struct {
	Message map[string]interface{} `json:"message"`
}

type PaginationCore struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Count    int    `json:"count"`
	NumPages int    `json:"num_pages"`
}

type CoursesListCore struct {
	Results    []CourseCore   `json:"results"`
	Pagination PaginationCore `json:"pagination"`
}

type CourseCore struct {
	ID               string                        `json:"id"`
	BlocksUrl        string                        `json:"blocks_url"`
	Effort           string                        `json:"effort"`
	EnrollmentStart  time.Time                     `json:"enrollment_start"`
	EnrollmentEnd    time.Time                     `json:"enrollment_end"`
	End              time.Time                     `json:"end"`
	Name             string                        `json:"name"`
	Number           string                        `json:"number"`
	Org              string                        `json:"org"`
	ShortDescription string                        `json:"short_description"`
	Start            time.Time                     `json:"start"`
	StartDisplay     string                        `json:"start_display"`
	StartType        string                        `json:"start_type"`
	Pacing           string                        `json:"pacing"`
	MobileAvailable  bool                          `json:"mobile_available"`
	Hidden           bool                          `json:"hidden"`
	InvitationOnly   bool                          `json:"invitation_only"`
	CourseID         string                        `json:"course_id"`
	Overview         string                        `json:"overview"`
	Media            *CourseApiMediaCollectionCore `json:"media"`
}

func (ht *CourseHTTP) FromCore(course *CourseCore) {
	ht.ID = course.ID
	ht.BlocksURL = course.BlocksUrl
	ht.Effort = course.Effort
	ht.EnrollmentStart = course.EnrollmentStart.String()
	ht.EnrollmentEnd = course.EnrollmentEnd.String()
	ht.Name = course.Name
	ht.Number = course.Number
	ht.Org = course.Org
	ht.ShortDescription = course.ShortDescription
	ht.Start = course.Start.String()
	ht.StartDisplay = course.StartDisplay
	ht.StartType = course.StartType
	ht.Pacing = course.Pacing
	ht.MobileAvailable = course.MobileAvailable
	ht.Hidden = course.Hidden
	ht.InvitationOnly = course.InvitationOnly
	ht.CourseID = course.CourseID
	ht.End = course.End.String()
	ht.Overview = &course.Overview
	ht.Media = &CourseAPIMediaCollectionHTTP{
		BannerImage: &AbsoluteMediaHTTP{},
		CourseImage: &MediaHTTP{},
		CourseVideo: &MediaHTTP{},
		Image:       &ImageHTTP{},
	}
	ht.Media.FromCore(course.Media)
}

func (ht *CourseHTTP) ToCore() *CourseCore {
	mediaCore := &CourseApiMediaCollectionCore{}
	mediaCore = ht.Media.ToCore()
	timeEnrollmentStart, _ := time.Parse("2006-Jan-02", ht.EnrollmentStart)
	timeEnrollmentEnd, _ := time.Parse("2006-Jan-02", ht.EnrollmentEnd)
	timeStart, _ := time.Parse("2006-Jan-02", ht.Start)
	timeEnd, _ := time.Parse("2006-Jan-02", ht.End)
	return &CourseCore{
		ID:               ht.ID,
		BlocksUrl:        ht.BlocksURL,
		Effort:           ht.Effort,
		EnrollmentStart:  timeEnrollmentStart,
		EnrollmentEnd:    timeEnrollmentEnd,
		Name:             ht.Name,
		Number:           ht.Number,
		Org:              ht.Org,
		ShortDescription: ht.ShortDescription,
		Start:            timeStart,
		StartDisplay:     ht.StartDisplay,
		StartType:        ht.StartType,
		Pacing:           ht.Pacing,
		MobileAvailable:  ht.MobileAvailable,
		Hidden:           ht.Hidden,
		InvitationOnly:   ht.InvitationOnly,
		CourseID:         ht.CourseID,
		End:              timeEnd,
		Media:            mediaCore,
		Overview:         *ht.Overview,
	}
}

func FromCoursesCore(courses []CourseCore) (coursesHttp []*CourseHTTP) {
	for _, courseCore := range courses {
		tmpCourseHttp := CourseHTTP{
			Media: &CourseAPIMediaCollectionHTTP{
				BannerImage: &AbsoluteMediaHTTP{},
				CourseImage: &MediaHTTP{},
				CourseVideo: &MediaHTTP{},
				Image:       &ImageHTTP{},
			},
		}
		tmpCourseHttp.FromCore(&courseCore)
		coursesHttp = append(coursesHttp, &tmpCourseHttp)
	}
	return coursesHttp
}
