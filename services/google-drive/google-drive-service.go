package driveservice

import (
	"context"
	"io"
	"log"
	"mime"
	"os"

	"google.golang.org/api/drive/v3"
)

// DriveDatastore manage all API action
type DriveDatastore struct {
	driveService *drive.Service
}

var (
	ctx = context.Background()
)

// NewDrive ...
func NewDrive() *DriveDatastore {
	currentSession := &DriveDatastore{nil}

	driveService, err := drive.NewService(ctx)
	if err != nil {
		log.Fatalf("Unable to retrieve GOOGLE_APPLICATION_CREDENTIALS %v", err)
	}

	currentSession.driveService = driveService
	log.Println("Connected to Google Drive")

	return currentSession
}

// List function return all files
func List(dr *DriveDatastore) interface{} {
	files, _ := dr.driveService.Files.List().Do()
	return files
}

// Upload function upload file to drive
func Upload(dr *DriveDatastore, name, mimeType, parentID string, content io.Reader) (*drive.File, error) {
	f := &drive.File{
		MimeType: mimeType,
		Name:     name,
		// Parents:  []string{parentID},
	}

	result, err := dr.driveService.Files.Create(f).Media(content).Do()
	if err != nil {
		log.Println("Could not create file: " + err.Error())
		return nil, err
	}

	return result, nil
}

// Download function will return a file base on fileID
func Download(dr *DriveDatastore, fileID string) interface{} {
	res, err := dr.driveService.Files.Get(fileID).Download()
	if err != nil {
		log.Println("Could not donwnload file: " + err.Error())
	}

	// Get file extension
	fileExtension, err := mime.ExtensionsByType(res.Header.Get("Content-Type"))
	if err != nil {
		log.Println("Could not get file extension: " + err.Error())
	}

	// Create empty file with extension
	outFile, err := os.Create("uname" + fileExtension[0])
	if err != nil {
		log.Println("Could not create file: " + err.Error())
	}
	defer outFile.Close()

	// Copy content to file that is created
	_, err = io.Copy(outFile, res.Body)
	if err != nil {
		log.Println("Could not copy content to file: " + err.Error())
	}

	return "uname" + fileExtension[0]
}

// Delete function will delete a file base on fileID
func Delete(dr *DriveDatastore, fileID string) error {
	err := dr.driveService.Files.Delete(fileID).Do()
	return err
}
