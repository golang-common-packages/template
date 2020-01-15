package onedriveservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/golang-common-packages/template/services"
)

// OneDriveService manage all drive action
type OneDriveService struct {
	Authentication *model.OneDrive
	URL            string
}

// NewOneDrive function return a new onedrive client based on singleton pattern
func NewOneDrive(config *model.Service) *OneDriveService {
	currentSession := &OneDriveService{nil, ""}

	oneDriveAuth := &model.OneDrive{
		AccessToken:  config.Database.OneDrive.AccessToken,
		RefreshToken: config.Database.OneDrive.RefreshToken,
	}

	currentSession.Authentication = oneDriveAuth
	currentSession.URL = config.Database.OneDrive.URL
	log.Println("Connected to OneDrive")

	return currentSession
}

// Search ...
func Search(od *OneDriveService, fileModel *model.FileModel) ([]byte, error) {
	deletePatch := fmt.Sprintf(od.URL+"/me/drive/search(q='%s')", filepath.Base(fileModel.Query))

	request := putRequest(deletePatch, od.Authentication.AccessToken, fileModel.Content)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	return body, nil
}

// List ...
func List(od *OneDriveService, fileModel *model.FileModel) ([]string, error) {
	childrenURL := od.URL + "/me/drive/root/children"

	if fileModel.Path != "" {
		parentFolderItem, err := itemByPath(od.Authentication.AccessToken, od.URL, fileModel.Path)
		if err != nil {
			return []string{}, err
		}
		childrenURL = fmt.Sprintf(od.URL+"/me/drive/items/%s/children", parentFolderItem.ID)
	}

	childItems, err := listItemsAsStruct(od.Authentication.AccessToken, childrenURL)
	if err != nil {
		return []string{}, err
	}

	items := []string{}
	for _, item := range childItems.Value {
		items = append(items, item.Name)
	}

	return items, nil
}

// Upload ...
func Upload(od *OneDriveService, fileModel *model.FileModel) (bool, error) {
	uploadFileURL := fmt.Sprintf(od.URL+"/me/drive/root:/%s:/content", filepath.Base(fileModel.Path))

	if filepath.Dir(fileModel.Path) != "." {
		parentFolderItem, err := itemByPath(od.Authentication.AccessToken, od.URL, filepath.ToSlash(filepath.Dir(fileModel.Path)))
		if err != nil {
			return false, err
		}
		uploadFileURL = fmt.Sprintf(od.URL+"/me/drive/items/%s:/%s:/content", parentFolderItem.ID, filepath.Base(fileModel.Path))
	}

	request := putRequest(uploadFileURL, od.Authentication.AccessToken, fileModel.Content)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return false, err
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return false, fmt.Errorf("Uploading the content was not successful. It returned the status code: %v", response.StatusCode)
	}

	return true, nil
}

// Download ...
func Download(od *OneDriveService, fileModel *model.FileModel) (string, error) {
	url := fmt.Sprintf(od.URL+"/me/drive/root:/%s:/content", fileModel.Path)
	request := getRequest(url, od.Authentication.AccessToken)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	return string(body), nil
}

// Delete ...
func Delete(od *OneDriveService, fileModel *model.FileModel) error {
	deletePatch := fmt.Sprintf(od.URL+"/me/drive/items/:%s:", filepath.Base(fileModel.Path))

	request := deleteRequest(deletePatch, od.Authentication.AccessToken)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return fmt.Errorf("Delete the content was not successful. It returned the status code: %v", response.StatusCode)
	}
	return nil
}

// Move ...
func Move(od *OneDriveService, fileModel *model.FileModel) ([]byte, []error) {
	var mvErrors []error
	createFolderPatch := fmt.Sprintf(od.URL+"/me/drive/items/%s", filepath.Base(fileModel.SourcesID))

	moveFolderInfo := &model.MoveOneDriveItem{
		Name:            fileModel.Name,
		ParentReference: model.ID{"{new-parent-folder-id}"},
	}

	// Convert struct to io.Reader
	requestByte, _ := json.Marshal(moveFolderInfo)
	requestReader := bytes.NewReader(requestByte)

	request := patchRequest(createFolderPatch, od.Authentication.AccessToken, requestReader)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		relocationError := fmt.Errorf("Error when move file: %v", err)
		mvErrors = append(mvErrors, relocationError)
		return nil, mvErrors
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	return body, nil
}

// CreateFolder ...
func CreateFolder(od *OneDriveService, fileModel *model.FileModel) ([]byte, error) {
	createFolderPatch := fmt.Sprintf(od.URL + "/me/drive/root/children")

	createFolderInfo := &model.CreateOneDriveFolder{
		Name:                           fileModel.Name,
		MicrosoftGraphConflictBehavior: "rename",
	}

	// Convert struct to io.Reader
	requestByte, _ := json.Marshal(createFolderInfo)
	requestReader := bytes.NewReader(requestByte)

	request := postRequest(createFolderPatch, od.Authentication.AccessToken, requestReader)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	return body, nil
}

// OneDrive util
// itemByPath ...
func itemByPath(accessToken string, url, path string) (model.OneDriveItem, error) {
	client := &http.Client{}
	itemByPathURL := fmt.Sprintf(url+"/me/drive/root:/%s", path)

	request := getRequest(itemByPathURL, accessToken)
	response, err := client.Do(request)
	if err != nil {
		return model.OneDriveItem{}, err
	}

	return unmarshallItemResponse(response)
}

// listItemsAsStruct ...
func listItemsAsStruct(accessToken string, url string) (model.ListOneDriveItem, error) {
	client := &http.Client{}
	request := getRequest(url, accessToken)
	response, err := client.Do(request)
	if err != nil {
		return model.ListOneDriveItem{}, err
	}
	return unmarshallListResponse(response)
}

// unmarshallListResponse ...
func unmarshallListResponse(response *http.Response) (model.ListOneDriveItem, error) {
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	var unmarshalledResponse model.ListOneDriveItem

	err := json.Unmarshal(body, &unmarshalledResponse)
	if err != nil {
		return model.ListOneDriveItem{}, err
	}
	return unmarshalledResponse, nil
}

// unmarshallItemResponse ...
func unmarshallItemResponse(response *http.Response) (model.OneDriveItem, error) {
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	var unmarshalledResponse model.OneDriveItem

	err := json.Unmarshal(body, &unmarshalledResponse)
	if err != nil {
		return model.OneDriveItem{}, err
	}
	return unmarshalledResponse, nil
}

// getRequest ...
func getRequest(url string, accessToken string) *http.Request {
	req, _ := http.NewRequest("GET", url, nil)
	bearerToken := "Bearer " + accessToken
	req.Header.Add("Authorization", bearerToken)
	return req
}

// postRequest ...
func postRequest(url string, accessToken string, content io.Reader) *http.Request {
	req, _ := http.NewRequest("POST", url, content)
	bearerToken := "Bearer " + accessToken
	req.Header.Add("Authorization", bearerToken)
	return req
}

// putRequest ...
func putRequest(url string, accessToken string, content io.Reader) *http.Request {
	req, _ := http.NewRequest("PUT", url, ioutil.NopCloser(content))
	bearerToken := "Bearer " + accessToken
	req.Header.Add("Authorization", bearerToken)
	return req
}

// patchRequest ...
func patchRequest(url string, accessToken string, content io.Reader) *http.Request {
	req, _ := http.NewRequest("PATCH", url, ioutil.NopCloser(content))
	bearerToken := "Bearer " + accessToken
	req.Header.Add("Authorization", bearerToken)
	return req
}

// deleteRequest ...
func deleteRequest(url string, accessToken string) *http.Request {
	req, _ := http.NewRequest("DELETE", url, nil)
	bearerToken := "Bearer " + accessToken
	req.Header.Add("Authorization", bearerToken)
	return req
}
