package cachestore

import (
	"time"

	"github.com/labstack/echo/v4"

	"github.com/golang-common-packages/template/model"

	"github.com/golang-common-packages/template/common/util/hash"
)

// Cachestore store function in cachestore package
type Cachestore interface {
	Middleware(hash hash.Storage) echo.MiddlewareFunc
	Get(key string) (string, error)
	Delete(key string) error
	Set(key string, value string, expire time.Duration) error
}

const (
	REDIS = iota
	// ARANGODB
)

// NewCachestore function for Factory Pattern
func NewCachestore(datastoreType int, config *model.Service) Cachestore {

	switch datastoreType {
	case REDIS:
		return NewRedisCachestore(config)
		// case ARANGODB:
		// 	return NewArangodbDatastore(config)
	}

	return nil
}
