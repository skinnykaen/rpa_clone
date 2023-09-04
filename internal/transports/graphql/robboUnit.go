package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.33

import (
	"context"
	"net/http"
	"strconv"

	"github.com/skinnykaen/rpa_clone/internal/consts"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"github.com/skinnykaen/rpa_clone/pkg/utils"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// CreateRobboUnit is the resolver for the CreateRobboUnit field.
func (r *mutationResolver) CreateRobboUnit(ctx context.Context, input models.NewRobboUnit) (*models.RobboUnitHTTP, error) {
	robboUnit, err := r.robboUnitService.CreateRobboUnit(models.RobboUnitCore{Name: input.Name, City: input.City})
	if err != nil {
		r.loggers.Err.Printf("%s", err.Error())
		return nil, &gqlerror.Error{
			Extensions: map[string]interface{}{
				"err": err,
			},
		}
	}
	robboUnitHttp := models.RobboUnitHTTP{}
	robboUnitHttp.FromCore(robboUnit)
	return &robboUnitHttp, nil
}

// UpdateRobboUnit is the resolver for the UpdateRobboUnit field.
func (r *mutationResolver) UpdateRobboUnit(ctx context.Context, input models.UpdateRobboUnit) (*models.RobboUnitHTTP, error) {
	atoi, err := strconv.Atoi(input.ID)
	if err != nil {
		r.loggers.Err.Printf("%s", err.Error())
		return nil, &gqlerror.Error{
			Extensions: map[string]interface{}{
				"err": utils.ResponseError{
					Code:    http.StatusBadRequest,
					Message: consts.ErrAtoi,
				},
			},
		}
	}
	robboUnit := models.RobboUnitCore{
		ID:   uint(atoi),
		Name: input.Name,
		City: input.City,
	}
	updatedRobboUnit, err := r.robboUnitService.UpdateRobboUnit(robboUnit)
	if err != nil {
		r.loggers.Err.Printf("%s", err.Error())
		return nil, &gqlerror.Error{
			Extensions: map[string]interface{}{
				"err": err,
			},
		}
	}
	robboUnitHttp := models.RobboUnitHTTP{}
	robboUnitHttp.FromCore(updatedRobboUnit)
	return &robboUnitHttp, nil
}

// DeleteRobboUnit is the resolver for the DeleteRobboUnit field.
func (r *mutationResolver) DeleteRobboUnit(ctx context.Context, id string) (*models.Response, error) {
	atoi, err := strconv.Atoi(id)
	if err != nil {
		r.loggers.Err.Printf("%s", err.Error())
		return nil, &gqlerror.Error{
			Extensions: map[string]interface{}{
				"err": utils.ResponseError{
					Code:    http.StatusBadRequest,
					Message: consts.ErrAtoi,
				},
			},
		}
	}
	if err := r.robboUnitService.DeleteRobboUnit(uint(atoi)); err != nil {
		r.loggers.Err.Printf("%s", err.Error())
		return nil, &gqlerror.Error{
			Extensions: map[string]interface{}{
				"err": err,
			},
		}
	}
	return &models.Response{Ok: true}, err
}

// GetRobboUnitByID is the resolver for the GetRobboUnitById field.
func (r *queryResolver) GetRobboUnitByID(ctx context.Context, id string) (*models.RobboUnitHTTP, error) {
	atoi, err := strconv.Atoi(id)
	if err != nil {
		r.loggers.Err.Printf("%s", err.Error())
		return nil, &gqlerror.Error{
			Extensions: map[string]interface{}{
				"err": utils.ResponseError{
					Code:    http.StatusBadRequest,
					Message: consts.ErrAtoi,
				},
			},
		}
	}
	robboUnit, err := r.robboUnitService.GetRobboUnitById(uint(atoi))
	if err != nil {
		r.loggers.Err.Printf("%s", err.Error())
		return nil, &gqlerror.Error{
			Extensions: map[string]interface{}{
				"err": err,
			},
		}
	}
	robboUnitHttp := models.RobboUnitHTTP{}
	robboUnitHttp.FromCore(robboUnit)
	return &robboUnitHttp, nil
}

// GetAllRobboUnitByAccessToken is the resolver for the GetAllRobboUnitByAccessToken field.
func (r *queryResolver) GetAllRobboUnitByAccessToken(ctx context.Context, page *int, pageSize *int) (*models.RobboUnitHTTPList, error) {
	clientId := ctx.Value(consts.KeyId).(uint)
	clientRole := ctx.Value(consts.KeyRole).(models.Role)

	robboUnits, countRows, err := r.robboUnitService.GetAllRobboUnits(page, pageSize, clientId, clientRole)
	if err != nil {
		r.loggers.Err.Printf("%s", err.Error())
		return nil, &gqlerror.Error{
			Extensions: map[string]interface{}{
				"err": err,
			},
		}
	}
	return &models.RobboUnitHTTPList{
		RobboUnits: models.FromRobboUnitsCore(robboUnits),
		CountRows:  int(countRows),
	}, nil
}

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *queryResolver) GetRobboUnitsByUnitAdminID(ctx context.Context, page *int, pageSize *int, unitAdminID string) (*models.RobboUnitHTTPList, error) {
	atoi, err := strconv.Atoi(unitAdminID)
	if err != nil {
		r.loggers.Err.Printf("%s", err.Error())
		return nil, &gqlerror.Error{
			Extensions: map[string]interface{}{
				"err": utils.ResponseError{
					Code:    http.StatusBadRequest,
					Message: consts.ErrAtoi,
				},
			},
		}
	}
	robboUnits, countRows, err := r.robboUnitService.GetAllRobboUnits(page, pageSize, uint(atoi), models.RoleUnitAdmin)
	if err != nil {
		r.loggers.Err.Printf("%s", err.Error())
		return nil, &gqlerror.Error{
			Extensions: map[string]interface{}{
				"err": err,
			},
		}
	}
	return &models.RobboUnitHTTPList{
		RobboUnits: models.FromRobboUnitsCore(robboUnits),
		CountRows:  int(countRows),
	}, nil
}