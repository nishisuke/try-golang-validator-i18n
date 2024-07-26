package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

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
	if !ok {
		log.Fatal("translator not found")
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := ja_translations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		log.Fatal(err)
	}

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

	err = validate.RegisterTranslation("datetime", trans, func(ut ut.Translator) error {
		return ut.Add("datetime", "{0} は {1} の形式ではありません", true)
	},
		func(ut ut.Translator, fe validator.FieldError) string {
			p := fe.Param()
			p = strings.Replace(p, "2006", "YYYY", 1)
			p = strings.Replace(p, "01", "MM", 1)
			p = strings.Replace(p, "02", "DD", 1)

			v, err := ut.T("datetime", fe.Field(), p)
			if err != nil {
				log.Println(err)
			}
			return v
		})
	if err != nil {
		log.Fatal(err)
	}

	// Level 0
	err = validate.StructCtx(ctx, val)
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

	for _, err := range ve {
		fmt.Println(err.Translate(trans))
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
