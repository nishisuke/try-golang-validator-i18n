package main

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/go-playground/locales/ja"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	ja_translations "github.com/go-playground/validator/v10/translations/ja"
)

type value struct {
	FamilyName string `validate:"required"`
	FirstName  string `validate:"min=10"`
	Color      string `validate:"iscolor"`
	Birthdate  string `validate:"datetime=2006-01-02"`
}

func main() {
	ctx := context.Background()
	val := value{
		FamilyName: "",
		FirstName:  "John",
		Color:      "red",
	}

	ja := ja.New()
	uni := ut.New(ja, ja)

	trans, ok := uni.GetTranslator("ja")
	_ = ok

	validate := validator.New(validator.WithRequiredStructEnabled())
	ja_translations.RegisterDefaultTranslations(validate, trans)

	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		dict := map[string]string{
			"FamilyName": "名字",
			"FirstName":  "名前",
			"Color":      "色",
			"Birthdate":  "生年月日",
		}
		if name, ok := dict[field.Name]; ok {
			return name
		}

		return field.Name
	})

	// Level 0
	err := validate.StructCtx(ctx, val)
	printError(err, trans)
}

func printError(err error, trans ut.Translator) {
	if err == nil {
		fmt.Println("No error")
		return
	}
	var ve validator.ValidationErrors

	if !errors.As(err, &ve) {
		fmt.Printf("Unknown error: %s\n", err)
		return
	}

	fmt.Println(ve.Translate(trans))

	for _, err := range ve {
		fmt.Println("Namespace", err.Namespace())
		fmt.Println("Field", err.Field())
		fmt.Println("StructNamespace", err.StructNamespace())
		fmt.Println("StructField", err.StructField())
		fmt.Println("Tag", err.Tag())
		fmt.Println("Type", err.Type())
		fmt.Println("Value", err.Value())
		fmt.Println("Param", err.Param())
		fmt.Println("----")
	}
}
