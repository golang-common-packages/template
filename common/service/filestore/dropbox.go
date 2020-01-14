package filestore

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"

	"github.com/golang-microservices/template/model"
)

// DropboxServices manage all API action
type DropboxServices struct {
	config *dropbox.Config
	client files.Client
}

/*
	@dropboxSession: Mapping between hash and config for singleton pattern
*/
var (
	dropboxSession = make(map[string]*DropboxServices)
)

const (
	chunkSize int64 = 1 << 24 // 1 << 24 = 16777216
)

// NewDropbox ...
func NewDropbox(config *model.Service) Filestore {
	hash := config.Hash()
	currentSession := dropboxSession[hash]
	if currentSession == nil {
		currentSession = &DropboxServices{nil, nil}
		config := dropbox.Config{
			Token:    config.Database.Dropbox.Token,
			LogLevel: dropbox.LogInfo,
		}
		currentSession.config = &config
		currentSession.client = files.New(config)
		dropboxSession[hash] = currentSession
		log.Println("Connected to Dropbox")
	}
	return currentSession
}

// Search ...
func (db *DropboxServices) Search(fileModel *model.FileModel) (interface{}, error) {
	return nil, nil
}

// Metadata ...
func (db *DropboxServices) Metadata(fileModel *model.FileModel) (interface{}, error) {
	arg := files.NewGetMetadataArg(fileModel.Path)

	res, err := db.client.GetMetadata(arg)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// List ...
func (db *DropboxServices) List(fileModel *model.FileModel) (interface{}, error) {
	var err error
	var entries []files.IsMetadata

	arg := files.NewListFolderArg(fileModel.Path)
	res, err := db.client.ListFolder(arg)
	if err != nil {
		return nil, err
	}

	if res != nil {
		entries = res.Entries
		for res.HasMore {
			arg := files.NewListFolderContinueArg(res.Cursor)

			res, err = db.client.ListFolderContinue(arg)
			if err != nil {
				return nil, err
			}

			entries = append(entries, res.Entries...)
		}
	}

	return entries, nil
}

// Upload ...
func (db *DropboxServices) Upload(fileModel *model.FileModel) (interface{}, error) {
	dst, err := validatePath(fileModel.Destination)
	if err != nil {
		return nil, err
	}

	contents, err := os.Open(fileModel.Source)
	if err != nil {
		return nil, err
	}
	defer contents.Close()

	contentsInfo, err := contents.Stat()
	if err != nil {
		return nil, err
	}

	commitInfo := files.NewCommitInfo(dst)
	commitInfo.Mode.Tag = "overwrite"

	// The Dropbox API only accepts timestamps in UTC with second precision.
	commitInfo.ClientModified = time.Now().UTC().Round(time.Second)

	// For large file transfer
	if contentsInfo.Size() > chunkSize {
		metaData, err := uploadChunked(db, contents, commitInfo, contentsInfo.Size())
		if err != nil {
			return nil, err
		}

		return metaData, nil
	}

	// For normal file size transfer
	metaData, err := db.client.Upload(commitInfo, contents)
	if err != nil {
		return nil, err
	}

	return metaData, nil
}

// Download ...
func (db *DropboxServices) Download(fileModel *model.FileModel) (interface{}, error) {
	src, err := validatePath(fileModel.Source)
	if err != nil {
		return nil, err
	}

	arg := files.NewDownloadArg(src)
	_, contents, err := db.client.Download(arg)
	if err != nil {
		return nil, err
	}
	defer contents.Close()

	return contents, nil
}

// Delete ...
func (db *DropboxServices) Delete(fileModel *model.FileModel) error {
	force := true // Force delete file
	var deletePaths []string

	// Validate remove paths before executing removal
	for i := range fileModel.Destinations {
		destination, err := validatePath(fileModel.Destinations[i])
		if err != nil {
			return err
		}

		pathMetaData, err := db.Metadata(fileModel)
		if err != nil {
			return err
		}

		if _, ok := pathMetaData.(*files.FileMetadata); !ok {
			folderArg := files.NewListFolderArg(destination)
			res, err := db.client.ListFolder(folderArg)
			if err != nil {
				return err
			}
			if len(res.Entries) != 0 && !force {
				return fmt.Errorf("rm: cannot remove ‘%s’: Directory not empty, use `--force` or `-f` to proceed", destination)
			}
		}
		deletePaths = append(deletePaths, destination)
	}

	// Execute removals
	for _, path := range deletePaths {
		arg := files.NewDeleteArg(path)

		if _, err := db.client.DeleteV2(arg); err != nil {
			return err
		}
	}

	return nil
}

// Move ...
func (db *DropboxServices) Move(fileModel *model.FileModel) (interface{}, []error) {
	var mvErrors []error
	var relocationArgs []*files.RelocationArg

	argsToMove := fileModel.Sources[0 : len(fileModel.Sources)+1]
	re := regexp.MustCompile("[^/]+$")
	for _, argument := range argsToMove {
		argumentFile := re.FindString(argument)
		lastCharDest := fileModel.Destination[len(fileModel.Destination)-1:]

		var err error
		var arg *files.RelocationArg

		if lastCharDest == "/" {
			arg, err = makeRelocationArg(argument, fileModel.Destination+argumentFile)
		} else {
			arg, err = makeRelocationArg(argument, fileModel.Destination)
		}

		if err != nil {
			relocationError := fmt.Errorf("Error validating move for %s to %s: %v", argument, fileModel.Destination, err)
			mvErrors = append(mvErrors, relocationError)
		} else {
			relocationArgs = append(relocationArgs, arg)
		}
	}

	for _, arg := range relocationArgs {
		if _, err := db.client.MoveV2(arg); err != nil {
			moveError := fmt.Errorf("Move error: %v", arg)
			mvErrors = append(mvErrors, moveError)
		}
	}

	return nil, mvErrors
}

// CreateFolder ...
func (db *DropboxServices) CreateFolder(fileModel *model.FileModel) (interface{}, error) {
	dst, err := validatePath(fileModel.Destination)
	if err != nil {
		return nil, err
	}

	arg := files.NewCreateFolderArg(dst)
	result, err := db.client.CreateFolderV2(arg)
	if err != nil {
		return nil, err
	}

	return result, nil
}

/// Drop util ///
// validatePath ...
func validatePath(p string) (path string, err error) {
	path = p
	if !strings.HasPrefix(path, "/") {
		path = fmt.Sprintf("/%s", path)
	}

	path = strings.TrimSuffix(path, "/")

	return
}

// makeRelocationArg ...
func makeRelocationArg(s string, d string) (arg *files.RelocationArg, err error) {
	src, err := validatePath(s)
	if err != nil {
		return
	}
	dst, err := validatePath(d)
	if err != nil {
		return
	}

	arg = files.NewRelocationArg(src, dst)

	return
}

// uploadChunked ...
func uploadChunked(db *DropboxServices, r io.Reader, commitInfo *files.CommitInfo, sizeTotal int64) (*files.FileMetadata, error) {
	res, err := db.client.UploadSessionStart(files.NewUploadSessionStartArg(),
		&io.LimitedReader{R: r, N: chunkSize})
	if err != nil {
		return nil, err
	}

	written := chunkSize

	for (sizeTotal - written) > chunkSize {
		cursor := files.NewUploadSessionCursor(res.SessionId, uint64(written))
		args := files.NewUploadSessionAppendArg(cursor)

		err = db.client.UploadSessionAppendV2(args, &io.LimitedReader{R: r, N: chunkSize})
		if err != nil {
			return nil, err
		}
		written += chunkSize
	}

	cursor := files.NewUploadSessionCursor(res.SessionId, uint64(written))
	args := files.NewUploadSessionFinishArg(cursor, commitInfo)

	metaData, err := db.client.UploadSessionFinish(args, r)
	if err != nil {
		return nil, err
	}

	return metaData, nil
}
