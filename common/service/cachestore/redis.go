package cachestore

import (
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"

	"github.com/golang-microservices/template/model"

	"github.com/golang-microservices/template/common/util/hash"
)

// RedisCacheStore manage all redis action
type RedisCacheStore struct {
	Client *redis.Client
	Prefix string
}

/*
	@sessionMapping: Mapping between model.Database and RedisCacheStore for singleton pattern
*/
var (
	sessionMapping = make(map[string]*RedisCacheStore)
)

// NewRedisCachestore function return redis client based on singleton pattern
func NewRedisCachestore(config *model.Service) Cachestore {
	hash := config.Hash()
	currentSession := sessionMapping[hash]
	if currentSession == nil {
		currentSession = &RedisCacheStore{nil, ""}
		client, err := currentSession.connect(config.Database.Redis)
		if err != nil {
			panic(err)
		} else {
			currentSession.Client = client
			currentSession.Prefix = config.Database.Redis.Prefix
			sessionMapping[hash] = currentSession
			log.Println("Connected to Redis Server")
		}
	}

	return currentSession
}

// connect private function establish redis connection
func (r *RedisCacheStore) connect(data model.Redis) (client *redis.Client, err error) {
	if r.Client == nil {
		client = redis.NewClient(&redis.Options{
			Addr:     data.Host,
			Password: data.Password,
			DB:       data.DB,
		})

		_, err := client.Ping().Result()
		if err != nil {
			log.Println("Fail to connect redis: ", err)
			return nil, err
		}
	} else {
		client = r.Client
		err = nil
	}
	return
}

// Middleware function will provide an echo middleware for Redis
func (r *RedisCacheStore) Middleware(hash hash.Storage) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get(echo.HeaderAuthorization)
			key := hash.SHA512(token)

			if val, err := r.Get(key); err != nil {
				log.Printf("Can not get accesstoken from redis in redis middleware: %s", err.Error())
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			} else if val == "" {
				return c.NoContent(http.StatusUnauthorized)
			}

			return next(c)
		}
	}
}

// Set function will set key and value
func (r *RedisCacheStore) Set(key string, value string, expire time.Duration) (err error) {
	err = r.Client.Set(r.Prefix+key, value, expire).Err()
	return
}

// Get function will get value based on the key provided
func (r *RedisCacheStore) Get(key string) (value string, err error) {
	value, err = r.Client.Get(r.Prefix + key).Result()
	return
}

// Delete function will delete value based on the key provided
func (r *RedisCacheStore) Delete(key string) (err error) {
	err = r.Client.Del(r.Prefix + key).Err()
	return
}

// Close function will close redis connection
func (r *RedisCacheStore) Close() {
	r.Close()
}
