package services

import (
	"encoding/json"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"github.com/skinnykaen/rpa_clone/pkg/utils"
	"net/http"
)

type CourseService interface {
	GetCoursesByUser(userId uint, userRole models.Role) (courses models.CoursesListCore, err error)
	GetCourseById(id string) (course models.CourseCore, err error)
}

type CourseServiceImpl struct {
	edxService EdxService
}

func (c CourseServiceImpl) GetCourseById(id string) (course models.CourseCore, err error) {
	body, err := c.edxService.GetCourseById(id)
	if err != nil {
		return models.CourseCore{}, err
	}

	if err := json.Unmarshal(body, &course); err != nil {
		return models.CourseCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return course, nil
}

func (c CourseServiceImpl) GetCoursesByUser(userId uint, userRole models.Role) (coursesList models.CoursesListCore, err error) {
	switch userRole.String() {
	case models.RoleSuperAdmin.String():
		body, err := c.edxService.GetAllCourses()
		if err != nil {
			return models.CoursesListCore{}, err
		}

		if err := json.Unmarshal(body, &coursesList); err != nil {
			return models.CoursesListCore{}, utils.ResponseError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
		return coursesList, nil
	}

	return models.CoursesListCore{}, nil
}
