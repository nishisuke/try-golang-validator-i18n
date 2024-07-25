package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-playground/locales/ja"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	ja_translations "github.com/go-playground/validator/v10/translations/ja"
)

type value struct {
	FamilyName string `validate:"required"`
	FirstName  string `validate:"min=10"`
	Color      string `validate:"iscolor"`
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
