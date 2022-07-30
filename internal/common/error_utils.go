package common

const ErrorGeneric = "GENERIC_ERROR"
const ErrorDuplicateEntity = "DUPLICATE_ENTITY"
const ErrorNoResult = "NO_RESULT"
const ErrorPasswordMismatch = "PASSWORD_MISMATCH"
const ErrorPayloadMissingRequiredFields = "PAYLOAD_MISSING_REQUIRED_FIELDS"
const ErrorUnparsablePayload = "UNPARSABLE_PAYLOAD"
const ErrorUserExists = "USER_ALREADY_EXISTS"
const ErrAuthTokenExpired = "AUTH_TOKEN_EXPIRED"

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
