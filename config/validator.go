package config

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type (
	SignUpRequest struct {
		Name        string `json:"name" validate:"required,min=2,max=100"`
		Email       string `json:"email" validate:"required,email"`
		CountryCode string `json:"country_code" validate:"required,oneof=BG US"`
		Password    string `json:"password" validate:"required,min=8,max=100"`
	}
	LoginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8,max=100"`
	}
	UpdateProfileRequest struct {
		Name        string `json:"name" validate:"required,min=2,max=100"`
		CountryCode string `json:"country_code" validate:"required,oneof=BG US"`
	}
	ChangePasswordRequest struct {
		OldPassword string `json:"old_password" validate:"required,min=8,max=100"`
		NewPassword string `json:"new_password" validate:"required,min=8,max=100"`
	}
	InsertExpenseTypeRequest struct {
		Name string `json:"name" validate:"required,min=2,max=100"`
	}
	InsertExpenseRequest struct {
		Amount      float64 `json:"amount" validate:"required,numeric"`
		Note        string  `json:"note" validate:"required,min=2,max=100"`
		TypeId      uint64  `json:"type_id" validate:"required,numeric"`
		WorkspaceId uint64  `json:"workspace_id" validate:"required,numeric"`
		CurrencyId  uint64  `json:"currency_id" validate:"required,numeric"`
	}
	GetExpenseByIdRequest struct {
		Id uint64 `params:"id" validate:"required,numeric"`
	}
	ErrorResponse struct {
		Error       bool
		FailedField string
		Tag         string
		Value       interface{}
	}

	XValidator struct {
		validator *validator.Validate
	}

	GlobalErrorHandlerResp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
)

var validate = validator.New()

var v = &XValidator{validator: validate}

func GlobalErrorHandler(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusBadRequest).JSON(GlobalErrorHandlerResp{
		Success: false,
		Message: err.Error(),
	})
}

func (v *XValidator) Validate(data interface{}) []ErrorResponse {
	var errors []ErrorResponse

	errs := validate.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			errors = append(errors, ErrorResponse{
				Error:       true,
				FailedField: err.Field(),
				Tag:         err.Tag(),
				Value:       err.Param(),
			})

		}
	}

	return errors
}

func ValidationResponse(data interface{}) error {
	if errs := v.Validate(data); len(errs) > 0 && errs[0].Error {
		errMsgs := make([]string, 0)

		for _, err := range errs {
			errMsgs = append(errMsgs, fmt.Sprintf(
				"[%s]: '%v' | Needs to implement '%s'",
				err.FailedField,
				err.Value,
				err.Tag,
			))
		}

		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: strings.Join(errMsgs, ", "),
		}
	}

	return nil
}

// register custom validations

func InitValidations() {
	RegisterCustomFieldValidations("password", "password", PasswordValidation)
}

func RegisterCustomFieldValidations(field string, tag string, fn validator.Func) {
	validate.RegisterValidation(tag, fn)
	validate.RegisterAlias(tag, field)
}

// custom validator functions

func PasswordValidation(fl validator.FieldLevel) bool {
	s := fl.Field().String()
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(s) >= 7 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}
