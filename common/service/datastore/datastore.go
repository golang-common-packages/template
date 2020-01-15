package datastore

import (
	"github.com/golang-common-packages/template/model"
	// "go.mongodb.org/mongo-driver/bson/primitive"
)

// Datastore interface for DAO function
type Datastore interface {
	// For user
	GetUsers(lastID, pageSize string) ([]model.User, error)
	GetUser(username string) (model.User, error)
	SaveUser(user model.User) error
	UpdateUser(user *model.User) error
	DeleteUser(userID interface{}) error
	ActiveUser(username string) error

	// For document
	GetDocuments(lastID, pageSize string) (documents []model.Document, err error)
	SaveDocuments(document model.Document) error
}

const (
	MONGODB = iota
	POSTGRES
)

// NewDatastore function for Factory Pattern
func NewDatastore(datastoreType int, config *model.Service) Datastore {

	switch datastoreType {
	case MONGODB:
		return NewMongoDBDatastore(config)
	case POSTGRES:
		return NewPostgresDatastore(config)
	}

	return nil
}
