package common

const (
	ErrorGeneric                      = "GENERIC_ERROR"
	ErrorInvalidInput                 = "INVALID_INPUT"
	ErrorDuplicateEntity              = "DUPLICATE_ENTITY"
	ErrorNotFound                     = "NOT_FOUND"
	ErrorForbidden                    = "FORBIDDEN"
	ErrorUnauthorised                 = "UNAUTHORISED"
	ErrorPasswordMismatch             = "PASSWORD_MISMATCH"
	ErrorPayloadMissingRequiredFields = "PAYLOAD_MISSING_REQUIRED_FIELDS"
	ErrorUnparsablePayload            = "UNPARSABLE_PAYLOAD"
	ErrorUserExists                   = "USER_ALREADY_EXISTS"
	ErrorAuthTokenExpired             = "AUTH_TOKEN_EXPIRED"
)

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
