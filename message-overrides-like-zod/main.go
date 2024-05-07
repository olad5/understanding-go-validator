package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/fr"
	ut "github.com/go-playground/universal-translator"

	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	fr_translations "github.com/go-playground/validator/v10/translations/fr"
)

// User contains user information

type User struct {
	FirstName      string     `json:"first_name" validate:"required"`
	LastName       string     `json:"last_name" validate:"required"`
	Age            uint8      `json:"age" validate:"gte=0,lte=130"`
	Email          string     `json:"email" validate:"required,email"`
	FavouriteColor string     `validate:"hexcolor|rgb|rgba"`
	Addresses      []*Address `validate:"required,dive,required"` // a person can have a home and cottage...
}

// Address houses a users address information
type Address struct {
	Street string `validate:"required"`
	City   string `validate:"required"`
	Planet string `validate:"required"`
	Phone  string `validate:"required"`
}

// use a single instance , it caches struct info
var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
)

func main() {
	en := en.New()
	fr := fr.New()
	uni = ut.New(en, en, fr)

	// this is usually know or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	transEn, _ := uni.GetTranslator("en")
	transFr, _ := uni.GetTranslator("fr")
	transEn.Add("{{first_name}}", "First Name", false)
	transFr.Add("{{first_name}}", "Prénom", false)
	transEn.Add("{{last_name}}", "Last Name", false)
	transFr.Add("{{last_name}}", "Nom de famille", false)
	// transEn.Add("{{email}}", "Last Name", false)
	transEn.Add("{{id}}", "Last Name", false)
	transEn.Add("{{customerId}}", "Last Name", false)
	validate = validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}
		println(" this is how name looks like: ", name)

		return name // what are these brackets used for
		// return "{{" + name + "}}" // what are these brackets used for
	})

	en_translations.RegisterDefaultTranslations(validate, transEn)
	fr_translations.RegisterDefaultTranslations(validate, transFr)
	translateOverrideHereNew(transEn) // yep you can specify your own in whatever locale you want!

	// build 'User' info, normally posted data etc...
	address := &Address{
		Street: "Eavesdown Docks",
		Planet: "Persphone",
		Phone:  "none",
		City:   "Unknown",
	}

	user := &User{
		FirstName:      "",
		LastName:       "",
		Age:            45,
		Email:          "",
		FavouriteColor: "#000",
		Addresses:      []*Address{address},
	}

	// returns InvalidValidationError for bad validation input, nil or ValidationErrors ( []FieldError )
	err := validate.Struct(user)
	if err != nil {

		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return
		}

		// errMsg := make(map[string]string)

		for _, err := range err.(validator.ValidationErrors) {

			/*
				fmt.Println("Namespace: " + err.Namespace())
				fmt.Println("Field: " + err.Field())
				fmt.Println("StructNamespace: " + err.StructNamespace()) // can differ when a custom TagNameFunc is registered or
				fmt.Println("StructField: " + err.StructField())         // by passing alt name to ReportError like below
				fmt.Println("Tag: " + err.Tag())
				fmt.Println("ActualTag: " + err.ActualTag())
				fmt.Println("Kind: ", err.Kind())
				fmt.Println("Type: ", err.Type())
				fmt.Println("Value: ", err.Value())
				fmt.Println("Param: " + err.Param())
				fmt.Println(err.Translate(transFr))
				fmt.Println()
			*/
			jsonKey := err.Field()
			// fieldName, _ := transFr.T(jsonKey)
			// message := strings.Replace(err.Translate(transFr), jsonKey, fieldName, -1)
			fmt.Println("jsonKey here bro: ", jsonKey)

			// jsonKey = jsonKey[2 : len(jsonKey)-2]
			// errMsg[jsonKey] = message
			// fmt.Println(jsonKey, ":", errMsg[jsonKey])
			fmt.Println(jsonKey + " must be defined")
		}

		// from here you can create your own error messages in whatever language you wish
		return
	}

	// TODO:TODO: since i can get the tag now, it is not time to make sure it is required
	// save user to database
}

func translateOverrideHereNew(trans ut.Translator) {
	// TODO:TODO: this works but I need to use the json tag and not the struct field
	// TODO:TODO: this is what i am looking for
	// source:https://github.com/syssam/go-playground-sample/blob/master/main.go
	validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} must have a value!", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		// t, _ := ut.T("required", fe.Field())
		fmt.Println("FieldError validator fe.Field here: 92", fe.Field())
		fmt.Println("FieldError validator fe.Tag here: 90", fe.Tag())
		fmt.Println("FieldError validator fe.ActualTag here: 90", fe.ActualTag())
		fmt.Println("FieldError validator fe.StructField here: 90", fe.StructField())
		fmt.Println("FieldError validator fe.Type here: 90", fe.Type())
		fmt.Println("FieldError validator fe.Kind here: 90", fe.Kind())

		// t, _ := ut.T("required", fe.Namespace())
		t, _ := ut.T("required", fe.Field())

		return t
	})

	type User struct {
		Username string `json:"dexter" validate:"required"`
	}

	var user User

	err := validate.Struct(user)
	if err != nil {

		errs := err.(validator.ValidationErrors)

		for _, e := range errs {
			// can translate each error one at a time.
			// TODO:TODO: this is where the error message matters
			fmt.Println("some translations things here: 223", e.Translate(trans))
		}
	}
}
