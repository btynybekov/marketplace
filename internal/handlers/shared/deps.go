package shared

import (
	"github.com/btynybekov/marketplace/config"
	"github.com/btynybekov/marketplace/internal/repository"
)

type Deps struct {
	Config config.EnvConfig
	Repos  repository.RepositorySet
}
