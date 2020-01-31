package model

import (
	"time"
)

// GetUserByQueryString store query string for user
type GetUserByQueryString struct {
	Username string `json:"username" query:"username"`
	Skip     string `json:"skip" query:"skip"`
	Limit    string `json:"limit" query:"limit"`
}

// User model work with database
type User struct {
	ID         interface{} `json:"id" bson:"_id" gorm:"unique;type:int; primary_key"`
	Name       string      `json:"name,omitempty" bson:"name,omitempty" validate:"required" gorm:"column:name"`
	Age        int         `json:"age,omitempty" bson:"age,omitempty" gorm:"column:age"`
	Username   string      `json:"username,omitempty" bson:"username,omitempty" gorm:"unique_index;column:username"`
	Password   *string     `json:"password,omitempty" bson:"password,omitempty"  gorm:"column:password"`
	Email      string      `json:"email,omitempty" bson:"email,omitempty" gorm:"unique_index;column:email"`
	IsActive   bool        `json:"isactive" bson:"isactive" gorm:"column:is_active"`
	Created    time.Time   `json:"created,omitempty" bson:"created,omitempty" gorm:"column:created"`
	Updated    time.Time   `json:"updated,omitempty" bson:"updated,omitempty" gorm:"column:updated"`
	Expiration time.Time   `json:"expiration,omitempty" bson:"expiration,omitempty" gorm:"column:expiration"`
}

// UserResult model for api response
type UserResult struct {
	ID       interface{} `json:"id" bson:"_id" gorm:"unique;type:int; primary_key"`
	Name     string      `json:"name,omitempty" bson:"name,omitempty" validate:"required" gorm:"column:name"`
	Age      int         `json:"age,omitempty" bson:"age,omitempty" gorm:"column:age"`
	Username string      `json:"username,omitempty" bson:"username,omitempty" gorm:"unique_index;column:username"`
	Email    string      `json:"email,omitempty" bson:"email,omitempty" gorm:"unique_index;column:email"`
	IsActive bool        `json:"isactive" bson:"isactive" gorm:"column:is_active"`
	Created  time.Time   `json:"created,omitempty" bson:"created,omitempty" gorm:"column:created"`
}

func (User *User) TableName() string {
	return "user"
}
