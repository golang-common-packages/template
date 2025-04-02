package bookHttpDelivery

import (
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/golang-common-packages/template/domain"
)

type BookValidator struct {
	validate *validator.Validate
}

func NewBookValidator() *BookValidator {
	v := validator.New()
	
	// Register custom validations
	v.RegisterValidation("notempty", func(fl validator.FieldLevel) bool {
		return strings.TrimSpace(fl.Field().String()) != ""
	})

	return &BookValidator{
		validate: v,
	}
}

// ValidateBook validates a single book
func (bv *BookValidator) ValidateBook(book *domain.Book) error {
	err := bv.validate.Struct(book)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}

		// Custom error messages
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "required":
				return &AppError{
					Message:    "Field " + err.Field() + " is required",
					StatusCode: http.StatusBadRequest,
				}
			case "notempty":
				return &AppError{
					Message:    "Field " + err.Field() + " cannot be empty",
					StatusCode: http.StatusBadRequest,
				}
			case "max":
				return &AppError{
					Message:    "Field " + err.Field() + " exceeds maximum length of " + err.Param(),
					StatusCode: http.StatusBadRequest,
				}
			case "excludesall":
				return &AppError{
					Message:    "Field " + err.Field() + " contains invalid characters",
					StatusCode: http.StatusBadRequest,
				}
			case "alphaunicode":
				return &AppError{
					Message:    "Field " + err.Field() + " must contain only letters",
					StatusCode: http.StatusBadRequest,
				}
			}
		}
	}
	return nil
}

// ValidateBooks validates a slice of books
func (bv *BookValidator) ValidateBooks(books *[]domain.Book) error {
	for _, book := range *books {
		if err := bv.ValidateBook(&book); err != nil {
			return err
		}
	}
	return nil
}
