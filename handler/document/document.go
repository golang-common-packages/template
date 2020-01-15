package document

import (
	"bytes"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/golang-microservices/cloud-storage"

	"github.com/golang-microservices/template/config"
	"github.com/golang-microservices/template/model"
)

// Handler manage all request and dependency
type Handler struct {
	*config.Environment
}

// New return a new Handler
func New(env *config.Environment) *Handler {
	return &Handler{env}
}

// Handler function will register all path to echo.Echo
func (h *Handler) Handler(e *echo.Group) {
	e.GET("/document", h.list(), h.JWT.Middleware(h.Config.Token.Accesstoken.PublicKey), h.Cache.Middleware(h.Hash), h.Monitor.Middleware())
	e.POST("/document", h.save(), h.JWT.Middleware(h.Config.Token.Accesstoken.PublicKey), h.Cache.Middleware(h.Hash), h.Monitor.Middleware())
	e.GET("/drive", h.files())
	// e.GET("/drive", driveController.List())
	e.POST("/drive", h.upload())
	e.DELETE("/drive", h.delete())
	e.GET("/drive/donwload", h.donwload())
}

// localhost:3000/api/v1/document?limit=3
// localhost:3000/api/v1/document?limit=3&lastid=5cee0e7af554bfbe838882c2
func (h *Handler) list() echo.HandlerFunc {
	return func(c echo.Context) error {
		documents, err := h.Database.GetDocuments(c.QueryParam("lastid"), c.QueryParam("limit"))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, documents)
	}
}

func (h *Handler) save() echo.HandlerFunc {
	return func(c echo.Context) error {
		validate := validator.New()
		request := model.Document{}

		// Bind request body to struct
		if err := c.Bind(&request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		// Validate request body struct
		if err := validate.Struct(request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err := h.Database.SaveDocuments(request); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.NoContent(http.StatusOK)
	}
}

func (h *Handler) files() echo.HandlerFunc {
	return func(c echo.Context) error {
		fileInfo := &cloudStorage.FileModel{Path: ""}
		files, err := h.Storage.List(fileInfo)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, files)
	}
}

func (h *Handler) upload() echo.HandlerFunc {
	return func(c echo.Context) error {
		name := c.FormValue("name")
		mimeType := c.FormValue("mimeType")
		parentID := c.FormValue("parentID")

		file, err := c.FormFile("file")
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		src, err := file.Open()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		defer src.Close()

		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, src); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		// This will return a value of type bytes.Reader which implements the io.Reader (and io.ReadSeeker) interface.
		// Don't worry about them not being the same "type". io.Reader is an interface and can be implemented by many different types.
		data := bytes.NewReader(buf.Bytes())

		fileInfo := &cloudStorage.FileModel{
			Name:     name,
			MimeType: mimeType,
			ParentID: parentID,
			Content:  data,
		}

		result, err := h.Storage.Upload(fileInfo)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, result)
	}
}

func (h *Handler) delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		fileInfo := &cloudStorage.FileModel{
			SourcesID: c.QueryParam("fileid"),
		}

		err := h.Storage.Delete(fileInfo)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}

		return c.NoContent(http.StatusOK)
	}
}

func (h *Handler) donwload() echo.HandlerFunc {
	return func(c echo.Context) error {
		fileInfo := &cloudStorage.FileModel{
			SourcesID: c.QueryParam("fileid"),
		}

		files, err := h.Storage.Download(fileInfo)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}

		return c.File(files.(string))
	}
}
