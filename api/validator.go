package api

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/jotabf/simplebank/util"
)

func isEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9.!#$%&’*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z]{2,})+$`)
	return re.MatchString(email)
}

func isAlphanumeric(s string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return re.MatchString(s)
}

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}
	return false
}

var validUser validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if user, ok := fieldLevel.Field().Interface().(string); ok {
		return isEmail(user) || isAlphanumeric(user)
	}
	return false
}
