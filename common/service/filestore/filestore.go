package filestore

import (
	"context"

	"github.com/golang-microservices/template/model"
)

// Filestore interface for API storage
type Filestore interface {
	Search(fileModel *model.FileModel) (interface{}, error)
	Metadata(fileModel *model.FileModel) (interface{}, error)
	List(fileModel *model.FileModel) (interface{}, error)
	Upload(fileModel *model.FileModel) (interface{}, error)
	Download(fileModel *model.FileModel) (interface{}, error)
	Delete(fileModel *model.FileModel) error
	Move(fileModel *model.FileModel) (interface{}, []error)
	CreateFolder(fileModel *model.FileModel) (interface{}, error)
}

var (
	ctx = context.Background()
)

/*
	@DRIVE: Google Drive service
	@DROPBOX: Dropbox service
	@ONEDRIVE: OneDrive service
	@SHAREPOINT: Sharepoint service
*/
const (
	DRIVE = iota
	DROPBOX
	ONEDRIVE
	SHAREPOINT
)

// NewFilestore function for Factory Pattern
func NewFilestore(storageType int, config *model.Service) Filestore {
	switch storageType {
	case DRIVE:
		return NewDrive(config)
	case DROPBOX:
		return NewDropbox(config)
	case ONEDRIVE:
		return NewOneDrive(config)
	case SHAREPOINT:
		return NewSharepoint(config)
	}

	return nil
}
