package datastore

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"

	"github.com/golang-microservices/template/model"
)

// PostgresDatastore manage all postgres action
type PostgresDatastore struct {
	Connection *gorm.DB
	Config     *model.Database
}

/*
	@sessionMapping: Mapping between model.Database and PostgresDatastore for singleton pattern
	@ctx: returns a non-nil, empty Context when it's unclear which Context to use or it is not yet available
*/
var (
	sessionMappingPostgres = make(map[string]*PostgresDatastore)
)

// NewPostgresDatastore function return a new postgres client based on singleton pattern
func NewPostgresDatastore(config *model.Service) Datastore {
	hash := config.Hash()
	currentSession := sessionMappingPostgres[hash]
	if currentSession == nil {
		currentSession = &PostgresDatastore{nil, nil}

		client, err := ConnectPostgres(config)
		if err != nil {
			log.Println("Can not ping to Postgres server: ", err)
			panic(err)
		} else {
			currentSession.Connection = client
			currentSession.Config = &config.Database
			sessionMappingPostgres[hash] = currentSession
			log.Println("Connected to Postgres Server")
		}
	}
	return currentSession
}

// ConnectPostgres function return a new postgres DB based on singleton pattern
func ConnectPostgres(config *model.Service) (*gorm.DB, error) {
	URI := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable", config.Database.Postgres.User, config.Database.Postgres.DBName, config.Database.Postgres.Password, config.Database.Postgres.Host, config.Database.Postgres.Port)

	connection, err := gorm.Open("postgres", URI)
	if err != nil {
		return nil, err
	}
	connection.AutoMigrate(&model.User{}, &model.Document{})

	return connection, nil
}

// SaveUser fucntion store user to table in Posgres
func (p *PostgresDatastore) SaveUser(user model.User) (err error) {
	lastRecord := model.User{}
	p.Connection.Last(&lastRecord)
	user.Created = time.Now()

	if lastRecord.ID != nil {
		user.ID = lastRecord.ID.(int64) + 1
	} else {
		user.ID = 0
	}

	// Check record before insert to table
	if status := p.Connection.NewRecord(user); status == false {
		return errors.New("Cannot insert a record when that record already exists (SaveUser function)")
	}

	if status := p.Connection.Create(&user); status.Error != nil {
		return status.Error
	}

	return nil
}

// ActiveUser fucntion update status of user to table in Posgres
func (p *PostgresDatastore) ActiveUser(username string) (err error) {
	user := model.User{}

	if status := p.Connection.Where("username = ? ", username).First(&user); status.Error != nil {
		return status.Error
	}

	user.IsActive = true
	user.Updated = time.Now()

	if status := p.Connection.Save(&user); status.Error != nil {
		return status.Error
	}

	return nil
}

// GetUser fucntion get user by username in Posgres
func (p *PostgresDatastore) GetUser(username string) (user model.User, err error) {

	if status := p.Connection.Where("username = ? ", username).Find(&user); status.Error != nil {
		return user, status.Error
	}

	return user, nil
}

// GetUsers fucntion get list user in Posgres
func (p *PostgresDatastore) GetUsers(lastID, pageSize string) (users []model.User, err error) {

	if lastID == "" && pageSize == "" {
		if status := p.Connection.Find(&users); status.Error != nil {
			return nil, status.Error
		}
	} else {
		if status := p.Connection.Offset(lastID).Order("id asc").Limit(pageSize).Find(&users); status.Error != nil {
			return nil, status.Error
		}
	}

	return users, nil
}

// UpdateUser fucntion update user in Posgres
func (p *PostgresDatastore) UpdateUser(user *model.User) (err error) {

	if status := p.Connection.Save(&user); status.Error != nil {
		return status.Error
	}

	return nil
}

// DeleteUser fucntion delete user in Posgres
func (p *PostgresDatastore) DeleteUser(userID interface{}) (err error) {
	user := model.User{}

	if status := p.Connection.Where("id = ? ", userID).Delete(&user); status.Error != nil {
		return status.Error
	}

	return nil
}

// GetDocuments fucntion get document in Posgres
func (p *PostgresDatastore) GetDocuments(lastID, pageSize string) (documents []model.Document, err error) {

	if lastID == "" {
		if status := p.Connection.Order("id asc").Limit(pageSize).Find(&documents); status.Error != nil {
			return nil, status.Error
		}
	} else {
		if status := p.Connection.Offset(lastID).Order("id asc").Limit(pageSize).Find(&documents); status.Error != nil {
			return nil, status.Error
		}
	}

	return documents, nil
}

// SaveDocuments fucntion create document in Posgres
func (p *PostgresDatastore) SaveDocuments(document model.Document) (err error) {
	lastRecord := model.Document{}
	p.Connection.Last(&lastRecord)
	document.Created = time.Now()

	if lastRecord.ID != nil {
		document.ID = lastRecord.ID.(int64) + 1
	} else {
		document.ID = 0
	}

	// Check record before insert to table
	if status := p.Connection.NewRecord(document); status == false {
		return errors.New("Cannot insert a record when that record already exists (SaveDocuments function)")
	}

	if status := p.Connection.Create(&document); status.Error != nil {
		return status.Error
	}

	return nil
}
