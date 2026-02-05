package validator

import (
	"github.com/dmitry/taskmanager/pkg/errors"
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func New() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

func (v *Validator) Validate(data interface{}) error {
	if err := v.validate.Struct(data); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		details := make([]errors.ErrorDetail, 0)

		for _, fieldError := range validationErrors {
			details = append(details, errors.ErrorDetail{
				Field:   fieldError.Field(),
				Message: getErrorMessage(fieldError),
			})
		}

		return errors.Validation("Ошибка валидации", details)
	}

	return nil
}

func getErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "Это поле обязательно для заполнения"
	case "email":
		return "Введите корректный email адрес"
	case "min":
		return "Значение слишком короткое"
	case "max":
		return "Значение слишком длинное"
	case "gt":
		return "Значение должно быть больше 0"
	case "uuid":
		return "Должен быть корректный UUID"
	case "oneof":
		return "Должно быть одно из допустимых значений"
	default:
		return "Некорректное значение"
	}
}
