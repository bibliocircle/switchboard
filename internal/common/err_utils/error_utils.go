package err_utils

type DetailedError struct {
	ErrorCode   string `json:"errorCode"`
	Description string `json:"description"`
}

func (e DetailedError) Error() string {
	return e.Description
}

func NewDetailedError(code string, description string) *DetailedError {
	return &DetailedError{
		ErrorCode:   code,
		Description: description,
	}
}

func WrapAsDetailedError(err error) *DetailedError {
	if err == nil {
		return nil
	}
	return &DetailedError{
		ErrorCode:   ErrorGeneric,
		Description: err.Error(),
	}
}
