package model

import val "github.com/go-ozzo/ozzo-validation"

func requiredIf(condition bool) val.RuleFunc {
	return func(value interface{}) error {
		if condition {
			return val.Validate(value, val.Required)
		}
		return nil
	}
}
