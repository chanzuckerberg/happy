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

func NewBadRequestError(message string) CustomError {
	return CustomError{
		code:    400,
		message: message,
	}
}

func NewUnauthorizedError(message string) CustomError {
	return CustomError{
		code:    401,
		message: message,
	}
}

func NewForbiddenError(message string) CustomError {
	return CustomError{
		code:    403,
		message: message,
	}
}

func NewNotFoundError(message string) CustomError {
	return CustomError{
		code:    404,
		message: message,
	}
}

func NewInternalServerError(message string) CustomError {
	return CustomError{
		code:    500,
		message: message,
	}
}
