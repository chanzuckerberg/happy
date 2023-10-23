package api

type CustomError struct {
	code    int
	message string
}

func (e CustomError) Error() string {
	return e.message
}
