package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"

	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

// use a single instance , it caches struct info
var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
)

func main() {
	en := en.New()
	uni = ut.New(en, en)

	// this is usually know or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	transEn, _ := uni.GetTranslator("en")
	validate = validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	en_translations.RegisterDefaultTranslations(validate, transEn)
	translateOverride(transEn) // yep you can specify your own in whatever locale you want!

	// type User struct {
	// 	Username string `json:"william" validate:"required"`
	// }
	// var user User
	type User struct {
		// FirstName      string `json:"first_name" validate:"required"`
		// LastName       string `json:"last_name" validate:"required"`
		// Age            uint8  `json:"age" validate:"gte=0,lte=130"`
		// Email          string `json:"email" validate:"required,email"`
		FavouriteColor string `validate:"hexcolor|rgb|rgba"`
		// Addresses      []*Address `validate:"required,dive,required"` // a person can have a home and cottage...
	}

	// user := User{FirstName: "william"}
	user := User{FavouriteColor: "#ffffbb"}
	err := validate.Struct(user)
	if err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			fmt.Println("some err detail: 11", errs[0].Tag())
			msg := errs[0].Translate(transEn)
			fmt.Println("the error i am getting the high level consumer: 111", msg)
			return
		}

		return
	}
}

func translateOverride(trans ut.Translator) {
	// TODO:TODO: add more tags
	requiredTag := "required"
	validate.RegisterTranslation(requiredTag, trans, func(ut ut.Translator) error {
		return ut.Add(requiredTag, "{0} must have a value!", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(requiredTag, fe.Field())

		return t
	})

	colorTag := "hexcolor|rgb|rgba"
	validate.RegisterTranslation(colorTag, trans, func(ut ut.Translator) error {
		return ut.Add(colorTag, "{0} must be a valid color", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(colorTag, fe.Field())

		return t
	})
}
