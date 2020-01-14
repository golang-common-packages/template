package dropboxservice

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
)

// DropboxFile manage all dropbox client API
type DropboxFile struct {
	config *dropbox.Config
	client files.Client
}

const (
	chunkSize int64 = 1 << 24 // 1 << 24 = 16777216
)

// NewFileService ...
func NewFileService(token string) *DropboxFile {
	currentSession := &DropboxFile{nil, nil}
	config := dropbox.Config{
		Token:    token,
		LogLevel: dropbox.LogInfo,
	}
	currentSession.client = files.New(config)
	currentSession.config = &config

	return currentSession
}

// GetFileMetadata ...
func GetFileMetadata(db *DropboxFile, path string) (files.IsMetadata, error) {
	arg := files.NewGetMetadataArg(path)

	res, err := db.client.GetMetadata(arg)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// ListFilesWithMetadata ...
func ListFilesWithMetadata(db *DropboxFile, args []string) ([]files.IsMetadata, error) {
	filePatch := ""
	var err error
	var entries []files.IsMetadata

	if len(args) > 0 {
		filePatch, err = validatePath(args[0])
		if err != nil {
			return nil, err
		}
	}

	arg := files.NewListFolderArg(filePatch)
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

// uploadChunked ...
func uploadChunked(db *DropboxFile, r io.Reader, commitInfo *files.CommitInfo, sizeTotal int64) (err error) {
	res, err := db.client.UploadSessionStart(files.NewUploadSessionStartArg(),
		&io.LimitedReader{R: r, N: chunkSize})
	if err != nil {
		return
	}

	written := chunkSize

	for (sizeTotal - written) > chunkSize {
		cursor := files.NewUploadSessionCursor(res.SessionId, uint64(written))
		args := files.NewUploadSessionAppendArg(cursor)

		err = db.client.UploadSessionAppendV2(args, &io.LimitedReader{R: r, N: chunkSize})
		if err != nil {
			return
		}
		written += chunkSize
	}

	cursor := files.NewUploadSessionCursor(res.SessionId, uint64(written))
	args := files.NewUploadSessionFinishArg(cursor, commitInfo)

	if _, err = db.client.UploadSessionFinish(args, r); err != nil {
		return
	}

	return
}

// upload ...
func upload(db *DropboxFile, args []string) (err error) {
	if len(args) == 0 || len(args) > 2 {
		return errors.New("`put` requires `src` and/or `dst` arguments")
	}

	src := args[0]

	// Default `dst` to the base segment of the source path; use the second argument if provided.
	dst := "/" + path.Base(src)
	if len(args) == 2 {
		dst, err = validatePath(args[1])
		if err != nil {
			return
		}
	}

	contents, err := os.Open(src)
	if err != nil {
		return
	}
	defer contents.Close()

	contentsInfo, err := contents.Stat()
	if err != nil {
		return
	}

	commitInfo := files.NewCommitInfo(dst)
	commitInfo.Mode.Tag = "overwrite"

	// The Dropbox API only accepts timestamps in UTC with second precision.
	commitInfo.ClientModified = time.Now().UTC().Round(time.Second)

	if contentsInfo.Size() > chunkSize {
		return uploadChunked(db, contents, commitInfo, contentsInfo.Size())
	}

	if _, err = db.client.Upload(commitInfo, contents); err != nil {
		return
	}

	return
}

// Download ...
func Download(db *DropboxFile, args []string) (err error) {
	if len(args) == 0 || len(args) > 2 {
		return errors.New("`get` requires `src` and/or `dst` arguments")
	}

	src, err := validatePath(args[0])
	if err != nil {
		return
	}

	// Default `dst` to the base segment of the source path; use the second argument if provided.
	dst := path.Base(src)
	if len(args) == 2 {
		dst = args[1]
	}
	// If `dst` is a directory, append the source filename.
	if f, err := os.Stat(dst); err == nil && f.IsDir() {
		dst = path.Join(dst, path.Base(src))
	}

	arg := files.NewDownloadArg(src)

	_, contents, err := db.client.Download(arg)
	if err != nil {
		return
	}
	defer contents.Close()

	f, err := os.Create(dst)
	if err != nil {
		return
	}
	defer f.Close()

	if _, err = io.Copy(f, contents); err != nil {
		return
	}

	return
}

// Delete ...
func Delete(db *DropboxFile, args []string) error {
	force := true // Force delete file
	if len(args) < 1 {
		return errors.New("rm: missing operand")
	}

	var deletePaths []string

	// Validate remove paths before executing removal
	for i := range args {
		path, err := validatePath(args[i])
		if err != nil {
			return err
		}

		pathMetaData, err := GetFileMetadata(db, path)
		if err != nil {
			return err
		}

		if _, ok := pathMetaData.(*files.FileMetadata); !ok {
			folderArg := files.NewListFolderArg(path)
			res, err := db.client.ListFolder(folderArg)
			if err != nil {
				return err
			}
			if len(res.Entries) != 0 && !force {
				return fmt.Errorf("rm: cannot remove ‘%s’: Directory not empty, use `--force` or `-f` to proceed", path)
			}
		}
		deletePaths = append(deletePaths, path)
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
func Move(db *DropboxFile, args []string) error {
	var destination string
	var argsToMove []string

	if len(args) > 2 {
		destination = args[len(args)-1]
		argsToMove = args[0 : len(args)-1]
	} else if len(args) == 2 {
		destination = args[1]
		argsToMove = append(argsToMove, args[0])
	} else {
		return fmt.Errorf("mv command requires a source and a destination")
	}

	var mvErrors []error
	var relocationArgs []*files.RelocationArg

	re := regexp.MustCompile("[^/]+$")
	for _, argument := range argsToMove {

		argumentFile := re.FindString(argument)
		lastCharDest := destination[len(destination)-1:]

		var err error
		var arg *files.RelocationArg

		if lastCharDest == "/" {
			arg, err = makeRelocationArg(argument, destination+argumentFile)
		} else {
			arg, err = makeRelocationArg(argument, destination)
		}

		if err != nil {
			relocationError := fmt.Errorf("Error validating move for %s to %s: %v", argument, destination, err)
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

	for _, mvError := range mvErrors {
		fmt.Fprintf(os.Stderr, "%v\n", mvError)
	}

	return nil
}

// CreateFolder ...
func CreateFolder(db *DropboxFile, args []string) (err error) {
	if len(args) != 1 {
		return errors.New("`mkdir` requires a `directory` argument")
	}

	dst, err := validatePath(args[0])
	if err != nil {
		return
	}

	arg := files.NewCreateFolderArg(dst)

	if _, err = db.client.CreateFolderV2(arg); err != nil {
		return
	}

	return
}
