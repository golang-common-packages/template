package model

import (
	"time"
)

// Document store a document's infomation
type Document struct {
	ID      interface{} `json:"id" bson:"_id,omitempty" gorm:"unique;type:int; primary_key" `
	Name    string      `json:"name,omitempty" bson:"name,omitempty" validate:"required" gorm:"column:name"`
	Author  string      `json:"author,omitempty" bson:"author,omitempty" validate:"required" gorm:"column:author"`
	Created time.Time   `json:"created,omitempty" bson:"created,omitempty" gorm:"column:created"`
	Updated time.Time   `json:"updated,omitempty" bson:"updated,omitempty" gorm:"column:updated"`
}

// TableName return table name
func (document *Document) TableName() string {
	return "document"
}
