package filestore

import (
	"io"
	"log"
	"mime"
	"os"

	"google.golang.org/api/drive/v3"

	"github.com/golang-microservices/template/model"
)

// DriveServices manage all drive action
type DriveServices struct {
	driveService *drive.Service
}

/*
	@driveSession: Mapping between hash and DriveServices for singleton pattern
*/
var (
	driveSession = make(map[string]*DriveServices)
)

// NewDrive function return a new mongo client based on singleton pattern
func NewDrive(config *model.Service) Filestore {
	hash := config.Hash()
	currentSession := driveSession[hash]
	if currentSession == nil {
		currentSession = &DriveServices{nil}

		driveService, err := drive.NewService(ctx)
		if err != nil {
			log.Fatalf("Unable to retrieve GOOGLE_APPLICATION_CREDENTIALS %v", err)
		}

		currentSession.driveService = driveService
		driveSession[hash] = currentSession
		log.Println("Connected to Google Drive")
	}

	return currentSession
}

// Search ...
// Search ...
func (dr *DriveServices) Search(fileModel *model.FileModel) (interface{}, error) {
	return nil, nil
}

// Metadata
func (dr *DriveServices) Metadata(fileModel *model.FileModel) (interface{}, error) {
	return nil, nil
}

// List function return all files
func (dr *DriveServices) List(fileModel *model.FileModel) (interface{}, error) {
	files, err := dr.driveService.Files.List().Do()
	return files, err
}

// Upload function upload file to drive
func (dr *DriveServices) Upload(fileModel *model.FileModel) (interface{}, error) {
	f := &drive.File{
		MimeType: fileModel.MimeType,
		Name:     fileModel.Name,
		// Parents:  []string{parentID},
	}

	result, err := dr.driveService.Files.Create(f).Media(fileModel.Content).Do()
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Download function will return a file base on fileID
func (dr *DriveServices) Download(fileModel *model.FileModel) (interface{}, error) {
	res, err := dr.driveService.Files.Get(fileModel.SourcesID).Download()
	if err != nil {
		return nil, err
	}

	// Get file extension
	fileExtension, err := mime.ExtensionsByType(res.Header.Get("Content-Type"))
	if err != nil {
		log.Println("Could not get file extension: " + err.Error())
	}

	// Create empty file with extension
	outFile, err := os.Create("uname" + fileExtension[0])
	if err != nil {
		return nil, err
	}
	defer outFile.Close()

	// Copy content to file that is created
	_, err = io.Copy(outFile, res.Body)
	if err != nil {
		log.Println("Could not copy content to file: " + err.Error())
	}

	return "uname" + fileExtension[0], nil
}

// Delete function will delete a file base on fileID
func (dr *DriveServices) Delete(fileModel *model.FileModel) error {
	err := dr.driveService.Files.Delete(fileModel.SourcesID).Do()
	return err
}

// Move function will move a file base on 'Sources' and 'Destination'
func (dr *DriveServices) Move(fileModel *model.FileModel) (interface{}, []error) {
	return nil, nil
}

// CreateFolder function will create a folder base on 'Destination'
func (dr *DriveServices) CreateFolder(fileModel *model.FileModel) (interface{}, error) {
	return nil, nil
}
