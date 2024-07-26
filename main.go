package main

import (
	"context"
	"errors"
	"log"
	"reflect"
	"strings"

	"github.com/go-playground/locales/ja"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	ja_translations "github.com/go-playground/validator/v10/translations/ja"
)

type user struct {
	FamilyName string `validate:"required"`
	FirstName  string `validate:"min=10"`
	Birthdate  string `validate:"datetime=2006-01-02"`
}

const langJA = "ja"

func main() {
	ctx := context.Background()
	val := user{
		FamilyName: "",
		FirstName:  "John",
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	ja := ja.New()
	uni := ut.New(ja, ja)
	err := registerTranslator(validate, uni, langJA)
	if err != nil {
		log.Fatal(err)
	}

	if err := validate.StructCtx(ctx, val); err != nil {
		printError(err, uni, langJA)
	}
}

func printError(err error, uni *ut.UniversalTranslator, lang string) {
	var ve validator.ValidationErrors

	if !errors.As(err, &ve) {
		log.Printf("Unknown error: %s\n", err)
		return
	}
	trans, ok := uni.GetTranslator(lang)
	if !ok {
		log.Println(ve)
	}

	for _, err := range ve {
		log.Println(err.Translate(trans))
	}
}

func registerTranslator(validate *validator.Validate, uni *ut.UniversalTranslator, lang string) error {
	trans, ok := uni.GetTranslator(lang)
	if !ok {
		log.Println("translator not found")
		return nil
	}

	// level1
	err := ja_translations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		return err
	}

	// level2
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

	// level3
	err = validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0}を入力してください", true)
	},
		func(ut ut.Translator, fe validator.FieldError) string {
			v, err := ut.T("required", fe.Field())
			if err != nil {
				log.Println(err)
			}
			return v
		})
	if err != nil {
		return err
	}

	// level4
	err = validate.RegisterTranslation("datetime", trans, func(ut ut.Translator) error {
		return ut.Add("datetime", "{0}は{1}の形式で入力してください", true)
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
		return err
	}
	return nil
}
