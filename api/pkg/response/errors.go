package response

type CustomError struct {
	code    int
	message string
}

func (e CustomError) Error() string {
	return e.message
}

func (e CustomError) GetCode() int {
	return e.code
}

func NewForbiddenError(message string) CustomError {
	return CustomError{
		code:    403,
		message: message,
	}
}

func NewUnauthorizedError(message string) CustomError {
	return CustomError{
		code:    401,
		message: message,
	}
}

func NewInternalServerError(message string) CustomError {
	return CustomError{
		code:    500,
		message: message,
	}
}
