package http

import (
	"github.com/skinnykaen/rpa_clone/internal/services"
	"github.com/skinnykaen/rpa_clone/pkg/logger"
)

type Handlers struct {
	ProjectHandler ProjectHandler
}

func SetupHandlers(
	loggers logger.Loggers,
	projectService services.ProjectService,
) Handlers {
	return Handlers{
		ProjectHandler: &ProjectHandlerImpl{
			loggers:        loggers,
			projectService: projectService,
		},
	}
}
