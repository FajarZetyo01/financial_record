package helpers

import (
	"database/sql"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type Validation struct {
	db *sql.DB
}

func NewValidation(db *sql.DB) *Validation {
	return &Validation{db: db}
}

func initValidation(validation *Validation) (*validator.Validate, ut.Translator) {
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	validate := validator.New()
	en_translations.RegisterDefaultTranslations(validate, trans)

	//CUSTOM LABEL FIELD
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		labelName := field.Tag.Get("label")
		if labelName == "" {
			return field.Name
		}
		return labelName
	})

	//CUSTOM TRANSLATE REQUIRED ALL
	validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} tidak boleh kosong", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	}) //REQUIRED ambil dari entities "struct"

	//CUSTOM TRANSLATE REQUIRED EMAIL
	validate.RegisterTranslation("email", trans, func(ut ut.Translator) error {
		return ut.Add("email", "{0} Email harus valid!", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email", fe.Field())
		return t
	})

	//CUSTOM MINIMAL PASSWORD
	validate.RegisterTranslation("min", trans, func(ut ut.Translator) error {
		return ut.Add("min", "{0} minimal {1} karakter", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("min", fe.Field(), fe.Param())
		return t
	})

	//CUSTOM CONFIRM PASSWORD EQFIELD
	validate.RegisterTranslation("eqfield", trans, func(ut ut.Translator) error {
		return ut.Add("eqfield", "{0} harus sama dengan {1} sebelumnya", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("eqfield", fe.Field(), fe.Param())
		return t
	})

	//REGISTER ISUNIQUE
	validate.RegisterValidation("isunique", func(fl validator.FieldLevel) bool {
		param := fl.Param()
		splitParam := strings.Split(param, "-")
		tableName := splitParam[0]
		fieldName := splitParam[1]
		fieldValue := fl.Field().String()
		query := "SELECT " + fieldName + " FROM " + tableName + " WHERE " + fieldName + " = ?"
		row := validation.db.QueryRow(query, fieldValue)
		var result string
		err := row.Scan(&result)
		return err == sql.ErrNoRows
	})
	//CUSTOM TEMPLATE ISUNIQUE
	validate.RegisterTranslation("isunique", trans, func(ut ut.Translator) error {
		return ut.Add("isunique", "{0} sudah terdaftar!", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("isunique", fe.Field())
		return t
	})
	return validate, trans
}

func (validation *Validation) ValidateStruct(s interface{}) interface{} {
	validate, trans := initValidation(validation)
	var validationError = make(map[string]interface{})

	if err := validate.Struct(s); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			validationError[e.StructField()] = e.Translate(trans)
		}
	}
	if len(validationError) > 0 {
		return validationError
	}
	return nil
}
