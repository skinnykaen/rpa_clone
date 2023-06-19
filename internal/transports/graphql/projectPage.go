package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.33

import (
	"context"
	"fmt"
	"github.com/skinnykaen/rpa_clone/internal/consts"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"strconv"
)

// CreateProjectPage is the resolver for the CreateProjectPage field.
func (r *mutationResolver) CreateProjectPage(ctx context.Context) (*models.ProjectPageHTTP, error) {
	newProjectPage, err := r.projectPageService.CreateProjectPage(ctx.Value(consts.KeyId).(uint))
	if err != nil {
		r.loggers.Err.Printf("%s", err.Error())
		return nil, err
	}
	projectPageHttp := models.ProjectPageHTTP{}
	projectPageHttp.FromCore(newProjectPage)
	return &projectPageHttp, nil
}

// UpdateProjectPage is the resolver for the UpdateProjectPage field.
func (r *mutationResolver) UpdateProjectPage(ctx context.Context, input models.UpdateProjectPage) (*models.ProjectPageHTTP, error) {
	panic(fmt.Errorf("not implemented: UpdateProjectPage - UpdateProjectPage"))
}

// DeleteProjectPage is the resolver for the DeleteProjectPage field.
func (r *mutationResolver) DeleteProjectPage(ctx context.Context, id string) (*models.Response, error) {
	atoi, err := strconv.Atoi(id)
	if err != nil {
		r.loggers.Err.Printf("%s", err.Error())
		return nil, err
	}
	if r.projectPageService.DeleteProjectPage(uint(atoi), ctx.Value(consts.KeyId).(uint)); err != nil {
		r.loggers.Err.Printf("%s", err.Error())
		return nil, err
	}
	return &models.Response{Ok: true}, nil
}

// GetProjectPageByID is the resolver for the GetProjectPageById field.
func (r *queryResolver) GetProjectPageByID(ctx context.Context, id string) (*models.ProjectPageHTTP, error) {
	atoi, err := strconv.Atoi(id)
	if err != nil {
		r.loggers.Err.Printf("%s", err.Error())
		return nil, err
	}
	project, err := r.projectPageService.GetProjectPageById(uint(atoi), ctx.Value(consts.KeyId).(uint))
	if err != nil {
		r.loggers.Err.Printf("%s", err.Error())
		return nil, err
	}
	projectPageHttp := models.ProjectPageHTTP{}
	projectPageHttp.FromCore(project)
	return &projectPageHttp, nil
}

// GetAllProjectPagesByAuthorID is the resolver for the GetAllProjectPagesByAuthorId field.
func (r *queryResolver) GetAllProjectPagesByAuthorID(ctx context.Context, id *string, page *int, pageSize *int) (*models.ProjectPageHTTPList, error) {
	panic(fmt.Errorf("not implemented: GetAllProjectPagesByAuthorID - GetAllProjectPagesByAuthorId"))
}

// GetAllProjectPagesByAccessToken is the resolver for the GetAllProjectPagesByAccessToken field.
func (r *queryResolver) GetAllProjectPagesByAccessToken(ctx context.Context, page *int, pageSize *int) (*models.ProjectPageHTTPList, error) {
	panic(fmt.Errorf("not implemented: GetAllProjectPagesByAccessToken - GetAllProjectPagesByAccessToken"))
}
