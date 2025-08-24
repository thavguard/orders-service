package service

import (
	"orders/src/db/repositories"
	"orders/src/mycache"
)

type Service struct {
	Repo    *repositories.DBRepository
	MyCache *mycache.RedisService
}
