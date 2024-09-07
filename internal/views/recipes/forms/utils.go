package forms

func isRequiredStyling(v bool) map[string]bool {
	return map[string]bool{
		"after:content-['_*']": v,
		"after:text-red-500":   v,
	}
}
