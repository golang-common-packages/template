package sharepointservice

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/api"
	strategy "github.com/koltyakov/gosip/auth/saml"

	"github.com/golang-microservices/template/services"
)

// SharepointService manage all sharepoint action
type SharepointService struct {
	Auth              *strategy.AuthCnfg
	Client            *gosip.SPClient
	SharepointService *api.SP
}

/*
	@sharepoinSession: Mapping between hash and SharepointService for singleton pattern
*/
var (
	headers = map[string]string{
		"Accept":          "application/json;odata=minimalmetadata",
		"Accept-Language": "de-DE,de;q=0.9",
	}
	config = &api.RequestConfig{Headers: headers}
)

// NewSharepoint function return a new sharepoint service based on singleton pattern
func NewSharepoint(config *model.Service) *SharepointService {
	currentSession := &SharepointService{nil, nil, nil}

	auth := &strategy.AuthCnfg{
		SiteURL:  config.Database.Sharepoint.SiteURL,
		Username: config.Database.Sharepoint.Username,
		Password: config.Database.Sharepoint.Password,
	}

	client := &gosip.SPClient{
		AuthCnfg: auth,
	}

	SharepointService := api.NewSP(client)

	currentSession.Auth = auth
	currentSession.Client = client
	currentSession.SharepointService = SharepointService
	log.Println("Connected to Sharepoint Server")

	return currentSession
}

// List function return all files with id and title
func List(sp *SharepointService, fileModel *model.FileModel) ([]byte, error) {
	// Assumes you have Custom list created
	endpoint := sp.Client.AuthCnfg.GetSiteURL() + "/_api/web/lists/getByTitle('Custom')/items"
	req, err := http.NewRequest(
		"POST",
		endpoint,
		bytes.NewBuffer([]byte(`{"__metadata":{"type":"SP.Data.CustomListItem"},"Title":"Test"}`)),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json;odata=verbose")
	req.Header.Set("Content-Type", "application/json;odata=verbose")
	resp, err := sp.Client.Execute(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Upload function upload file to sharepoint
func Upload(sp *SharepointService, fileModel *model.FileModel) ([]byte, error) {
	endpoint := fmt.Sprintf(
		"%s/_api/web/getFolderByServerRelativeUrl('%s')/files/add(overwrite=true,url='%s')",
		sp.Client.AuthCnfg.GetSiteURL(),
		url.QueryEscape(fileModel.ParentID),
		url.QueryEscape(fileModel.Name),
	)

	byteContent, err := ioutil.ReadAll(fileModel.Content)
	if err != nil {
		log.Fatalf("Unable to read %v", err)
	}

	req, err := http.NewRequest(
		"POST",
		endpoint,
		bytes.NewBuffer(byteContent),
	)
	if err != nil {
		log.Fatalf("Unable to POST data %v", err)
	}

	req.Header.Set("Accept", "application/json;odata=verbose")

	resp, err := sp.Client.Execute(req)
	if err != nil {
		log.Fatalf("Unable to execute %v", err)
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Unable to read body %v", err)
	}

	return result, nil
}

// Download function will return a file. Now, fileID is File URI that can be host relevant (e.g. `/sites/site/lib/folder/file.txt`)
func Download(sp *SharepointService, fileModel *model.FileModel) ([]byte, error) {
	endpoint := fmt.Sprintf(
		"%s/_api/Web/GetFileByServerRelativeUrl(@FileServerRelativeUrl)/$value?@FileServerRelativeUrl='%s'",
		sp.Client.AuthCnfg.GetSiteURL(),
		url.QueryEscape(fileModel.SourcesID),
	)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.TransferEncoding = []string{"null"}

	resp, err := sp.Client.Execute(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Delete function will delete a file with fileID
func Delete(sp *SharepointService, fileModel *model.FileModel) error {
	_, err := sp.SharepointService.Conf(config).Web().GetFile(fileModel.SourcesID).Recycle()
	return err
}
