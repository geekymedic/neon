package v9

func ValidateStruct(obj interface{}) error {
	var validate = defaultValidator{}
	return validate.ValidateStruct(obj)
}
