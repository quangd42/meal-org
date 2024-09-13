package forms

// HACK: This is not applied to InputText. If change styling here, remember to change it in InputText as well.
func isRequiredStyling(v bool) map[string]bool {
	return map[string]bool{
		"after:content-['_*']": v,
		"after:text-red-500":   v,
	}
}
